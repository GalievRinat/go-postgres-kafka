# go-postgres-kafka
Тестовое задание

Сервис принимает GET запросы по адресу /api/newmessage в JSON-виде:
```json
{
    "topic": "topic_name",
    "title": "title",
    "comment": "comment"
}
```
Сохраняет их в БД postgresql и передает в Apache Kafka
