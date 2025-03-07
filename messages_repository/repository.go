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
	//fmt.Println(psqlInfo)
	messagesRepo.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Ошибка открытия БД: ", err)
		return err
	}
	//var version int
	//err = messagesRepo.DB.QueryRow("SHOW server_version_num").Scan(&version)
	//if err != nil {
	//fmt.Println("Ошибка открытия БД: ", err)
	//return err
	//}
	//fmt.Printf("БД [%s] на хосте [%s:%d] успешно подключена, версия postgres %d\n", dbName, host, port, version)
	return nil
}

func (messagesRepo *MessagesRepository) Add(message model.Message) (int64, error) {
	fmt.Printf("Добавление сообщения в БД\n")
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
		fmt.Printf("Сообщение c ID [%d] не найдено: \n", message.ID)
		return errors.New("message with id not found")
	}

	return nil
}

func (messagesRepo *MessagesRepository) TotalCount() (int64, error) {
	var totalCount int64
	err := messagesRepo.DB.QueryRow("SELECT COUNT(*) FROM all_messages").Scan(&totalCount)
	if err != nil {
		fmt.Println("Ошибка запроса количества сообщений: ", err)
		return -1, err
	}
	return totalCount, nil
}

func (messagesRepo *MessagesRepository) SendCount() (int64, error) {
	var sendCount int64
	err := messagesRepo.DB.QueryRow("SELECT COUNT(*) FROM all_messages WHERE sendtokafka = True").Scan(&sendCount)
	if err != nil {
		fmt.Println("Ошибка запроса количества отправленных сообщений: ", err)
		return -1, err
	}
	return sendCount, nil
}

func (messagesRepo *MessagesRepository) GetMessageByID(id int64) (model.Message, error) {
	var message model.Message
	row := messagesRepo.DB.QueryRow("SELECT id, timestamp, topic, title, comment, sendtokafka FROM all_messages WHERE id=$1", id)
	err := row.Scan(&message.ID, &message.Timestamp, &message.Topic, &message.Title, &message.Comment, &message.SendToKafka)
	if err != nil {
		fmt.Println("Ошибка запроса сообщения: ", err)
		return model.Message{}, err
	}
	return message, nil
}

func (messagesRepo *MessagesRepository) GetAllMessages(count int64) ([]model.Message, error) {
	var messages []model.Message
	rows, err := messagesRepo.DB.Query("SELECT id, timestamp, topic, title, comment, sendtokafka FROM all_messages ORDER BY timestamp LIMIT $1", count)

	if err != nil {
		fmt.Println("Ошибка запроса сообщений: ", err)
		return []model.Message{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var message model.Message
		err := rows.Scan(&message.ID, &message.Timestamp, &message.Topic, &message.Title, &message.Comment, &message.SendToKafka)
		if err != nil {
			fmt.Println(err)
			return []model.Message{}, err
		}
		messages = append(messages, message)
	}
	if rows.Err() != nil {
		fmt.Println("Ошибка запроса сообщений: ", err)
		return []model.Message{}, rows.Err()
	}
	return messages, nil
}

func (messagesRepo *MessagesRepository) GetUnsendMessages() ([]model.Message, error) {
	var messages []model.Message
	rows, err := messagesRepo.DB.Query("SELECT id, timestamp, topic, title, comment, sendtokafka FROM all_messages WHERE sendtokafka IS NOT TRUE ORDER BY timestamp")

	if err != nil {
		fmt.Println("Ошибка запроса неотправленных сообщений: ", err)
		return []model.Message{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var message model.Message
		err := rows.Scan(&message.ID, &message.Timestamp, &message.Topic, &message.Title, &message.Comment, &message.SendToKafka)
		if err != nil {
			fmt.Println(err)
			return []model.Message{}, err
		}
		messages = append(messages, message)
	}
	if rows.Err() != nil {
		fmt.Println("Ошибка запроса неотправленных сообщений: ", err)
		return []model.Message{}, rows.Err()
	}
	return messages, nil
}
