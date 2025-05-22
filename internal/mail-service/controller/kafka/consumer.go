package kafka

import (
	"context"
	"log"
	"retarget/pkg/entity/notice"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
)

type Consumer struct {
	processor *goka.Processor
}

func NewConsumer(brokers []string, group string, topic goka.Stream) *Consumer {
	cb := func(ctx goka.Context, msg interface{}) {
		if event, ok := msg.(*notice.NoticeEvent); ok {
			log.Printf("\n[KAFKA] NEW MESSAGE\n"+
				"Topic: %s\n"+
				"Partition: %d\n"+
				"Offset: %d\n"+
				"Key: %s\n"+
				"User ID: %d\n"+
				"Type: %d\n"+
				"Message: %s\n"+
				"-------------------------",
				ctx.Topic(), ctx.Partition(), ctx.Offset(), ctx.Key(),
				event.UserID, event.Type, event.Message)
		} else {
			log.Printf("Received unknown message format: %v", msg)
		}
	}

	input := goka.Input(topic, new(codec.String), cb)

	processor, err := goka.NewProcessor(brokers,
		goka.DefineGroup(goka.Group(group),
			input,
		),
		goka.WithConsumerGroupBuilder(goka.DefaultConsumerGroupBuilder),
	)

	if err != nil {
		log.Fatalf("error creating processor: %v", err)
	}

	return &Consumer{
		processor: processor,
	}
}

func (c *Consumer) Run(ctx context.Context) {
	if err := c.processor.Run(ctx); err != nil {
		log.Printf("error running processor: %v", err)
	}
}
