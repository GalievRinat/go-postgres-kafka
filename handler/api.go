package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GalievRinat/go-postgres-kafka/model"
)

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

	id, err := handler.messagesRepo.Add(message)
	if err != nil || id < 0 {
		jsonError(w, "Ошибка добавления задачи в БД", err)
		return
	}

	message.ID = id

	err = handler.KafkaNewMessage(message)
	if err != nil {
		jsonError(w, "Ошибка записи сообщения в Kafka:", err)
		return
	}

	err = handler.messagesRepo.MarkSend(message)
	if err != nil {
		return
	}
	fmt.Println("Задача добавлена")
	w.WriteHeader(http.StatusOK)
	jsonError(w, "Задача добавлена", nil)
}
