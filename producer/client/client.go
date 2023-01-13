package client

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	log "github.com/sirupsen/logrus"
	"time"
)

type KafkaClient struct {
	w kafka.Writer
}

func NewKafkaClient(brokers []string, topic string) *KafkaClient {

	username := "alice"
	password := "alice-secret"
	partition := 0
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: username,
			Password: password,
		},
	}

	_, err := dialer.DialLeader(context.Background(), "tcp", brokers[0], topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	client := &KafkaClient{}

	w := kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
		Transport: &kafka.Transport{
			SASL: plain.Mechanism{
				Username: "alice",
				Password: "alice-secret",
			},
		},
	}
	client.w = w
	return client
}

func (client *KafkaClient) SendMessage(ctx context.Context, message []byte) error {
	err := client.w.WriteMessages(ctx, kafka.Message{
		Value: message,
	})

	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
