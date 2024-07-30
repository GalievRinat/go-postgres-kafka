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
http://195.133.1.93:5000/api/stats

Просмотр всех сообщений:
http://195.133.1.93:5000/api/getallmessages
http://195.133.1.93:5000/api/getallmessages?count=10

Просмотр сообщения по ID:
http://195.133.1.93:5000/api/getmessage?id=1

Ответ отправляется в формате JSON
```
{"OK":"Задача добавлена"}
```