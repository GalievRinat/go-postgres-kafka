package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GalievRinat/go-postgres-kafka/messages_repository"
	"github.com/GalievRinat/go-postgres-kafka/model"
	"github.com/segmentio/kafka-go"
)

type Handler struct {
	messagesRepo *messages_repository.MessagesRepository
	messageKafka *kafka.Conn
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

func (handler *Handler) ApiNewMessage(w http.ResponseWriter, r *http.Request) {
	var message model.Message
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message.Timestamp = time.Now()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	r_count, err := handler.messagesRepo.Add(message)
	if err != nil || r_count == 0 {
		fmt.Println(err)
		jsonError(w, "Ошибка добавления задачи в БД", err)
		return
	}

	jsonError(w, "Задача добавлена", nil)
	w.WriteHeader(http.StatusOK)
}
