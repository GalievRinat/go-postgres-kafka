package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	err = handler.SendMessageToKafka(message)
	if err != nil || id < 0 {
		fmt.Println("Ошибка отправки сообщения в kafka: ", err)
	}

	fmt.Println("Задача добавлена")
	w.WriteHeader(http.StatusOK)
	json_text, err := json.Marshal(map[string]string{"OK": "Задача добавлена"})
	if err != nil {
		fmt.Println("Ошибка генерации JSON для jsonError:", err)
		return
	}

	_, err = w.Write(json_text)
	if err != nil {
		fmt.Println("Ошибка записи JSON:", err)
		return
	}
}

func (handler *Handler) ApiStats(w http.ResponseWriter, r *http.Request) {
	totalCount, err := handler.messagesRepo.TotalCount()
	if err != nil || totalCount < 0 {
		jsonError(w, "Ошибка подсчета количества сообщений", err)
		return
	}
	fmt.Println("Всего сообщений: ", totalCount)

	sendCount, err := handler.messagesRepo.SendCount()
	if err != nil || sendCount < 0 {
		jsonError(w, "Ошибка подсчета количества отправленных сообщений", err)
		return
	}
	fmt.Println("Всего отправленных сообщений: ", sendCount)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json_text, err := json.Marshal(map[string]string{
		"total_messages": strconv.FormatInt(totalCount, 10),
		"send_messages":  strconv.FormatInt(sendCount, 10),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Ошибка генерации JSON для jsonError:", err)
		return
	}

	_, err = w.Write(json_text)
	if err != nil {
		fmt.Println("Ошибка записи JSON:", err)
		return
	}
}

func (handler *Handler) ApiGetMessage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		jsonError(w, "Неправильный ID", err)
		return
	}
	message, err := handler.messagesRepo.GetMessageByID(id)
	if err != nil {
		jsonError(w, "Ошибка получения сообщения", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json_text, err := json.Marshal(message)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Ошибка генерации JSON для jsonError:", err)
		return
	}

	_, err = w.Write(json_text)
	if err != nil {
		fmt.Println("Ошибка записи JSON:", err)
		return
	}
}

func (handler *Handler) ApiGetAllMessages(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.ParseInt(r.URL.Query().Get("count"), 10, 64)
	if err != nil {
		count = 100
	}
	messages, err := handler.messagesRepo.GetAllMessages(count)
	if err != nil {
		jsonError(w, "Ошибка получения сообщения", err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json_text, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Ошибка генерации JSON для jsonError:", err)
		return
	}

	_, err = w.Write(json_text)
	if err != nil {
		fmt.Println("Ошибка записи JSON:", err)
		return
	}
}
