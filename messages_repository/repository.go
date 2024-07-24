package messages_repository

import (
	"database/sql"
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
	return nil
}

func (messagesRepo *MessagesRepository) Add(message model.Message) (int64, error) {
	res, err := messagesRepo.DB.Exec("INSERT INTO all_messages (timestamp, topic, title, comment, sendtokafka) VALUES ($1, $2, $3, $4, $5)",
		message.Timestamp, message.Topic, message.Title, message.Comment, message.SendToKafka)
	if err != nil {
		fmt.Println("Ошибка добавления задачи (insert): ", err)
		return 0, err
	}
	r_count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return r_count, nil
}
