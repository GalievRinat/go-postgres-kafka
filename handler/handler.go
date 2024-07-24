package handler

import (
	"fmt"

	"github.com/GalievRinat/go-postgres-kafka/messages_repository"
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
