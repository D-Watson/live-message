package mq

import (
	"context"
	"fmt"
	"time"

	cf "github.com/D-Watson/live-safety/conf"
	"github.com/segmentio/kafka-go"
)

type Client struct {
	writer *kafka.Writer
	reader *kafka.Reader
	topic  string
}

func NewProducer(topic string) *Client {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cf.GlobalConfig.Kafka.Address...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireAll,
		AllowAutoTopicCreation: true,
		BatchTimeout:           10 * time.Millisecond,
		BatchSize:              100,
		Async:                  false,
	}

	return &Client{
		writer: writer,
		topic:  topic,
	}
}

func NewConsumer(topic, groupID string) *Client {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cf.GlobalConfig.Kafka.Address,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        2 * time.Second,
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second, // 自动提交偏移量
	})

	return &Client{
		reader: reader,
		topic:  topic,
	}
}

func (c *Client) Produce(ctx context.Context, value []byte) error {
	msg := kafka.Message{
		Value: value,
		Time:  time.Now(),
	}

	if err := c.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	return nil
}

func (c *Client) Consume(ctx context.Context) (kafka.Message, error) {
	msg, err := c.reader.ReadMessage(ctx)
	if err != nil {
		return kafka.Message{}, fmt.Errorf("failed to consume message: %w", err)
	}
	return msg, nil
}

func (c *Client) ConsumeWithHandler(ctx context.Context, handler func([]byte) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				return fmt.Errorf("consumption error: %w", err)
			}

			if err := handler(msg.Value); err != nil {
				return fmt.Errorf("handler error: %w", err)
			}
		}
	}
}
func (c *Client) Close() error {
	var errs []error

	if c.writer != nil {
		if err := c.writer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close writer: %w", err))
		}
	}

	if c.reader != nil {
		if err := c.reader.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close reader: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during close: %v", errs)
	}
	return nil
}
func (c *Client) GetStats() interface{} {
	if c.writer != nil {
		return c.writer.Stats()
	}
	if c.reader != nil {
		return c.reader.Stats()
	}
	return nil
}
