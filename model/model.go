package model

import "time"

type Message struct {
	ID          int64     `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Topic       string    `json:"topic"`
	Title       string    `json:"title"`
	Comment     string    `json:"comment"`
	SendToKafka bool      `json:"sendtokafka"`
}
