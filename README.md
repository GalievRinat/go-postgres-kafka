# go-postgres-kafka
Тестовое задание

Сервис принимает GET запросы по адресу /api/newmessage в формате JSON:
```
{
    "topic": "topic_one",
    "title": "Title",
    "comment": "Comment"
}
```

Просмотр статистики сообщений:
- /api/stats

Просмотр всех сообщений:
- /api/getallmessages
- /api/getallmessages?count=10

Просмотр сообщения по ID:
- /api/getmessage?id=1

Ответ отправляется в формате JSON
```
{"OK":"Задача добавлена"}
```