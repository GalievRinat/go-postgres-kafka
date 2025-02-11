package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/GalievRinat/go-postgres-kafka/handler"
	"github.com/go-chi/chi/v5"

	//gotdotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//err := gotdotenv.Load()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	port, err := strconv.Atoi(os.Getenv("GPK_DBPORT"))
	if err != nil {
		fmt.Println("Ошибка чтения порта БД:", err)
		return
	}
	sendInterval, err := strconv.Atoi(os.Getenv("KAFKA_SENDINTERVAL"))
	if err != nil {
		fmt.Println("Ошибка чтения интервала отправки в kafka:", err)
		return
	}

	host := os.Getenv("GPK_DBHOST")
	user := os.Getenv("GPK_DBUSER")
	password := os.Getenv("GPK_DBPASSWORD")
	dbName := os.Getenv("GPK_DBNAME")

	handler, err := handler.NewHandler(host, port, user, password, dbName)
	if err != nil {
		fmt.Println("Ошибка создания handler: ", err)
		return
	}
	defer handler.CloseHandler()
	go handler.SendTiker(sendInterval)

	r := chi.NewRouter()

	r.Get("/api/newmessage", handler.ApiNewMessage)
	r.Get("/api/stats", handler.ApiStats)
	r.Get("/api/getmessage", handler.ApiGetMessage)
	r.Get("/api/getallmessages", handler.ApiGetAllMessages)
	addr := fmt.Sprintf(":%s", os.Getenv("GPK_APIPORT"))
	fmt.Printf("Start web server on port [%s]\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		fmt.Printf("Start server error: %s", err.Error())
		return
	}
}
