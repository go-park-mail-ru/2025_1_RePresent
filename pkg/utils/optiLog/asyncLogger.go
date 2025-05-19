package optiLog

import (
	"context"

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
	logChan chan *LogMessage
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *zap.SugaredLogger
}

func NewAsyncLogger(logger *zap.SugaredLogger, bufferSize int) *AsyncLogger {
	ctx, cancel := context.WithCancel(context.Background())
	al := &AsyncLogger{
		logChan: make(chan *LogMessage, bufferSize),
		ctx:     ctx,
		cancel:  cancel,
		logger:  logger,
	}
	go al.processLogs()
	return al
}

func (al *AsyncLogger) processLogs() {
	for {
		select {
		case msg := <-al.logChan:
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
		case <-al.ctx.Done():
			close(al.logChan)
			return
		}
	}
}

func (al *AsyncLogger) Log(level zapcore.Level, requestID, message string, fields map[string]interface{}) {
	if cap(al.logChan) == 0 {
		return
	}
	al.logChan <- &LogMessage{
		Level:     level,
		RequestID: requestID,
		Message:   message,
		Fields:    fields,
	}
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
