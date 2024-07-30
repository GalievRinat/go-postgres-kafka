# go-postgres-kafka
Сервис принимает GET запросы по адресу /api/newmessage в формате JSON:
```
{
    "topic": "topic_one",
    "title": "Title",
    "comment": "Comment"
}
```

Просмотр статистики сообщений:
```
/api/stats
```

Просмотр всех сообщений:
```
/api/getallmessages
/api/getallmessages?count=10
```

Просмотр сообщения по ID:
```
/api/getmessage?id=1
```

Ответ отправляется в формате JSON
```
{"OK":"Задача добавлена"}
```

Для запуска необходимо задать следующие переменные окружения:
```
GPK_APIPORT=Порт для подключения к микросервису

GPK_DBHOST=АДРЕС POSGTRESQL
GPK_DBPORT=ПОРТ POSTGRESQL
GPK_DBNAME=Имя БД
GPK_DBUSER=Пользователь БД
GPK_DBPASSWORD=Пароль

KAFKA_ADDR=Адрес сервера Apache Kafka
KAFKA_PORT=Порт

KAFKA_SENDINTERVAL=Интервал повторной отправки сообщений в Kafka, секунды
```

Контейнер с микросервисом:
```
viral0249/go-postgres-kafka
```
