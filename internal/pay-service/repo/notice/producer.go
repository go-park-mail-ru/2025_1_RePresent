package notice

import (
	"encoding/json"
	"fmt"
	"time"

	"retarget/pkg/entity/notice"
	noticeType "retarget/pkg/entity/notice"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"go.uber.org/zap"
)

type NoticeRepositoryInterface interface {
	SendLowBalanceNotification(userID int, message string) error
	SendTopUpBalanceEvent(userID int, message string) error
	Close()
}

type NoticeRepository struct {
	emitter     *goka.Emitter
	noticeTopic goka.Stream
	logger      *zap.SugaredLogger
}

func retryNewEmitter(brokers []string, stream goka.Stream, logger *zap.SugaredLogger) (*goka.Emitter, error) {
	var emitter *goka.Emitter
	var err error
	delay := time.Second

	for i := 1; i <= 5; i++ {
		emitter, err = goka.NewEmitter(brokers, stream, new(codec.String))
		if err == nil {
			logger.Infof("Connected to Kafka on attempt %d", i)
			return emitter, nil
		}
		logger.Warnf("Attempt %d: failed to create emitter: %v", i, err)
		time.Sleep(delay)
		delay *= 2
	}

	return nil, fmt.Errorf("could not connect to Kafka after 5 attempts: %w", err)
}

func NewNoticeRepository(brokers []string, topic string, logger *zap.SugaredLogger) *NoticeRepository {
	stream := goka.Stream(topic)
	emitter, err := retryNewEmitter(brokers, stream, logger)
	if err != nil {
		logger.Errorw("failed to create goka emitter", "error", err)
		return nil
	}

	return &NoticeRepository{
		emitter:     emitter,
		noticeTopic: stream,
		logger:      logger,
	}
}

func (r *NoticeRepository) SendTopUpBalanceEvent(userID int, amount float64) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if r.emitter == nil {
		return fmt.Errorf("emitter not initialized")
	}

	event := notice.NoticeEvent{
		UserID: userID,
		Type:   noticeType.TopUpedBalance,
		Amount: amount,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		r.logger.Errorw("failed to marshal notice event", "event", event, "error", err)
		return err
	}

	key := fmt.Sprintf("%d", userID)

	err = r.emitter.EmitSync(key, string(payload))
	if err != nil {
		r.logger.Errorw("failed to emit notice event", "key", key, "payload", string(payload), "error", err)
		return err
	}

	r.logger.Infow("notice event sent", "key", key, "payload", string(payload))
	return nil
}

func (r *NoticeRepository) SendLowBalanceNotification(userID int) error {
	if userID <= 0 {
		return fmt.Errorf("invalid user ID: %d", userID)
	}

	if r.emitter == nil {
		return fmt.Errorf("emitter not initialized")
	}

	event := notice.NoticeEvent{
		UserID: userID,
		Type:   noticeType.LowBalance,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		r.logger.Errorw("failed to marshal notice event", "event", event, "error", err)
		return err
	}

	key := fmt.Sprintf("%d", userID)

	err = r.emitter.EmitSync(key, string(payload))
	if err != nil {
		r.logger.Errorw("failed to emit notice event", "key", key, "payload", string(payload), "error", err)
		return err
	}

	r.logger.Infow("notice event sent", "key", key, "payload", string(payload))
	return nil
}

func (r *NoticeRepository) Close() {
	if r.emitter != nil {
		if err := r.emitter.Finish(); err != nil {
			r.logger.Errorw("Failed to finish emitter: %v", err)
		}
	}
}
