package handler

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/GalievRinat/go-postgres-kafka/model"
	"github.com/segmentio/kafka-go"
)

func (handler *Handler) KafkaNewMessage(message model.Message) error {
	addr := os.Getenv("KAFKA_ADDR")
	port := os.Getenv("KAFKA_PORT")
	addr_full := fmt.Sprintf("%s:%s", addr, port)

	topic := message.Topic
	partition := 0
	fmt.Printf("Добавление сообщения в Kafka [%s]\n", addr_full)
	conn, err := kafka.DialLeader(context.Background(), "tcp", addr_full, topic, partition)
	if err != nil {
		fmt.Println("Ошибка вызова Kafka:", err)
		return err
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	fullTextMessage := fmt.Sprintf("Title: %s\nComment: %s", message.Title, message.Comment)
	kafkaMessage := kafka.Message{Value: []byte(fullTextMessage)}
	_, err = conn.WriteMessages(kafkaMessage)
	if err != nil {
		fmt.Println("Ошибка записи сообщения в Kafka:", err)
		return err
	}

	if err := conn.Close(); err != nil {
		fmt.Println("Ошибка закрытия подключения к Kafka:", err)
		return err
	}
	return nil
}
