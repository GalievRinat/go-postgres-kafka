package handler

import (
	"fmt"
	"time"

	"github.com/GalievRinat/go-postgres-kafka/messages_repository"
	"github.com/GalievRinat/go-postgres-kafka/model"
)

type Handler struct {
	messagesRepo *messages_repository.MessagesRepository
}

func NewHandler(host string, port int, user string, password string, dbName string) (*Handler, error) {
	handler := Handler{}
	handler.messagesRepo = &messages_repository.MessagesRepository{}
	err := handler.messagesRepo.CreateRepo(host, port, user, password, dbName)
	return &handler, err
}

func (handler *Handler) CloseHandler() {
	err := handler.messagesRepo.DB.Close()
	fmt.Println("Ошибка закрытия handler: ", err)
}

func (handler *Handler) SendMessageToKafka(message model.Message) error {
	err := handler.KafkaNewMessage(message)
	if err != nil {
		return err
	}

	err = handler.messagesRepo.MarkSend(message)
	if err != nil {
		return err
	}
	return nil
}

func (handler *Handler) SendUnsendMessages() error {
	unsendMessages, err := handler.messagesRepo.GetUnsendMessages()
	if err != nil {
		return err
	}
	if unsendMessages == nil {
		fmt.Println("Нет неотправленных сообщений")
		return nil
	}
	errCount := 0
	for _, message := range unsendMessages {
		err := handler.SendMessageToKafka(message)
		if err != nil {
			fmt.Println("Ошибка отправки в kafka: ", err)
			errCount++
		}
	}
	if errCount != 0 {
		fmt.Println("Количество ошибок отправки в kafka: ", errCount)
	}
	return nil
}

func (handler *Handler) SendTiker(sendInterval int) {
	for range time.NewTicker(time.Duration(sendInterval) * time.Second).C {
		handler.SendUnsendMessages()
	}
}
