package optiLog

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogMessage struct {
	Level     zapcore.Level
	RequestID string
	Message   string
	Fields    map[string]interface{}
}

type AsyncLogger struct {
	ctx        context.Context
	cancel     context.CancelFunc
	logger     *zap.SugaredLogger
	mu         sync.Mutex
	logChan    chan *LogMessage
	bufferSize int
	maxSize    int
}

func NewAsyncLogger(logger *zap.SugaredLogger, initialBufferSize, maxBufferSize int) *AsyncLogger {
	ctx, cancel := context.WithCancel(context.Background())
	al := &AsyncLogger{
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
		bufferSize: initialBufferSize,
		maxSize:    maxBufferSize,
		logChan:    make(chan *LogMessage, initialBufferSize),
	}
	go al.processLogs()
	return al
}

func (al *AsyncLogger) processLogs() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-al.logChan:
			if !ok {
				return
			}
			al.writeLog(msg)

		case <-ticker.C:
			// Периодически проверяем, не пора ли увеличить буфер
			al.resizeBufferIfNeeded()

		case <-al.ctx.Done():
			al.drainAndClose()
			return
		}
	}
}

func (al *AsyncLogger) writeLog(msg *LogMessage) {
	fields := []interface{}{}
	for k, v := range msg.Fields {
		fields = append(fields, k, v)
	}

	switch msg.Level {
	case zapcore.DebugLevel:
		al.logger.Debugw(msg.Message, fields...)
	case zapcore.InfoLevel:
		al.logger.Infow(msg.Message, fields...)
	case zapcore.WarnLevel:
		al.logger.Warnw(msg.Message, fields...)
	case zapcore.ErrorLevel:
		al.logger.Errorw(msg.Message, fields...)
	default:
		al.logger.Infow(msg.Message, fields...)
	}
}

func (al *AsyncLogger) resizeBufferIfNeeded() {
	al.mu.Lock()
	defer al.mu.Unlock()

	// Если буфер заполнен более чем на 80%, удваиваем его
	if cap(al.logChan) >= len(al.logChan)+1 && cap(al.logChan) < al.maxSize {
		newSize := cap(al.logChan) * 2
		if newSize > al.maxSize {
			newSize = al.maxSize
		}

		oldChan := al.logChan
		newChan := make(chan *LogMessage, newSize)

		// Перемещаем оставшиеся логи в новый канал
		for {
			select {
			case msg := <-oldChan:
				newChan <- msg
			default:
				al.logChan = newChan
				al.bufferSize = newSize
				close(oldChan)
				return
			}
		}
	}
}

func (al *AsyncLogger) drainAndClose() {
	al.mu.Lock()
	defer al.mu.Unlock()

	close(al.logChan)
	for msg := range al.logChan {
		al.writeLog(msg)
	}
}

func (al *AsyncLogger) Log(level zapcore.Level, requestID, message string, fields map[string]interface{}) {
	msg := &LogMessage{
		Level:     level,
		RequestID: requestID,
		Message:   message,
		Fields:    fields,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	select {
	case al.getLogChan() <- msg:
	case <-ctx.Done():
		// Пропускаем лог, если не успеваем
		return
	}
}

func (al *AsyncLogger) getLogChan() chan *LogMessage {
	al.mu.Lock()
	defer al.mu.Unlock()
	return al.logChan
}

func (al *AsyncLogger) Close() {
	al.cancel()
}

func MakeLogFields(requestID string, durationMs int64, extra ...map[string]interface{}) map[string]interface{} {
	fields := map[string]interface{}{
		"request_id":  requestID,
		"timeTakenMs": durationMs,
	}
	for _, e := range extra {
		for k, v := range e {
			fields[k] = v
		}
	}
	return fields
}
