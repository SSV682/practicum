package clientConsumer

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	"time"
)

type KafkaClient struct {
	r *kafka.Reader
}

func NewKafkaClient(brokers []string, topic string) *KafkaClient {
	client := &KafkaClient{}
	client.r = kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: "test-reader",
		Topic:   topic,
	})
	return client
}

func (client *KafkaClient) Read(ctx context.Context) ([]byte, error) {
	msg, err := client.r.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}

func (client *KafkaClient) ReadBatch(ctx context.Context, timeout time.Duration) ([][]byte, error) {
	cctx, cancel := context.WithCancel(ctx)
	go func(cancel context.CancelFunc) {
		select {
		case <-time.After(timeout):
			cancel()
		}
	}(cancel)

	response := [][]byte{}
	for {
		select {
		default:
			msg, err := client.Read(cctx)
			if err != nil && err != context.Canceled {
				return nil, err
			} else if err == context.Canceled {
				return response, nil
			}
			log.Info("read msg")
			response = append(response, msg)
		case <-ctx.Done():
			fmt.Println("canceled")
			return response, nil
		}
	}

}
