package messages_repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/GalievRinat/go-postgres-kafka/model"
	_ "github.com/lib/pq"
)

type MessagesRepository struct {
	DB *sql.DB
}

func (messagesRepo *MessagesRepository) CreateRepo(host string, port int, user string, password string, dbName string) error {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
	messagesRepo.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Ошибка открытия БД: ", err)
		return err
	}
	fmt.Printf("БД [%s] на хосте [%s:%d] успешно подключена\n", dbName, host, port)
	return nil
}

func (messagesRepo *MessagesRepository) Add(message model.Message) (int64, error) {
	var lastInsertId int64
	err := messagesRepo.DB.QueryRow("INSERT INTO all_messages (timestamp, topic, title, comment, sendtokafka) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		message.Timestamp, message.Topic, message.Title, message.Comment, message.SendToKafka).Scan(&lastInsertId)
	if err != nil {
		fmt.Println("Ошибка добавления задачи (insert): ", err)
		return -1, err
	}
	return lastInsertId, nil
}

func (messagesRepo *MessagesRepository) MarkSend(message model.Message) error {
	res, err := messagesRepo.DB.Exec("UPDATE all_messages SET sendtokafka = True WHERE id = $1", message.ID)
	if err != nil {
		fmt.Println("Ошибка пометки сообщения отправленным: ", err)
		return err
	}
	r_count, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Ошибка пометки сообщения отправленным: ", err)
		return err
	}
	if r_count == 0 {
		fmt.Printf("Сообщение c ID [%s] не найдено: \n", message.ID)
		return errors.New("message with id not found")
	}

	return nil
}
