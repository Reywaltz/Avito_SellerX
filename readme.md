# Тестовое задание на позицию стажера-бекендера

Цель задания – разработать чат-сервер, предоставляющий HTTP API для работы с чатами и сообщениями пользователя.

## Инструкция по запуску:
Для запуска приложения необходимо развернуть контейнеры через команду docker-compose, находясь в директории с проектом. 
```bash
git clone https://github.com/Reywaltz/Avito_SellerX
cd Avito_SellerX
docker-compose up
```
Для настройки сервиса доступны следующие переменные среды:
  * DEV - конфигурация режима вывода сообщений логгера (True - для удобного вывода сообщений в виде текста. False для сообщения в формате JSON)
 * CONN_DB - строка подключения к базе данных.

После запуска контейнеров пользователю доступна пустая база данных, структура которой описана в директории sql в файле 0001_schema.sql:
```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    username TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS chats (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    chat INTEGER REFERENCES chats(id) ON DELETE CASCADE,
    author INTEGER REFERENCES users(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS users_chats (
    chat_id INTEGER REFERENCES chats(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE
);
```

 Если необходимо проверить работоспособность приложения на тестовых данных, то в той же директории присутствует файл fill_data.sql с SQL-запросами на вставку. Для того, чтобы его применить необходимо подключиться к базе данных (Например, через pgcli) и запустить данный скрипт:
 ```bash
pgcli -h localhost -p 5433 avito avito_user

\i sql/fill_data.sql
 ```

## Основные API методы

Методы обрабатывают HTTP POST запросы c телом, содержащим все необходимые параметры в JSON.

### Добавить нового пользователя

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"username": "user_1"}' \
  http://localhost:9000/users/add
```

Ответ: `id` созданного пользователя или HTTP-код ошибки.

### Создать новый чат между пользователями

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"name": "chat_1", "users": ["<USER_ID_1>", "<USER_ID_2>"]}' \
  http://localhost:9000/chats/add
```

Ответ: `id` созданного чата или HTTP-код ошибки.

Количество пользователей не ограничено.

### Отправить сообщение в чат от лица пользователя

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"chat": "<CHAT_ID>", "author": "<USER_ID>", "text": "hi"}' \
  http://localhost:9000/messages/add
```

Ответ: `id` созданного сообщения или HTTP-код ошибки.

### Получить список чатов конкретного пользователя

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"user": "<USER_ID>"}' \
  http://localhost:9000/chats/get
```

Ответ: cписок всех чатов со всеми полями, отсортированный по времени создания последнего сообщения в чате (от позднего к раннему). Или HTTP-код ошибки.

### Получить список сообщений в конкретном чате

Запрос:

```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"chat": "<CHAT_ID>"}' \
  http://localhost:9000/messages/get
```

Ответ: список всех сообщений чата со всеми полями, отсортированный по времени создания сообщения (от раннего к позднему). Или HTTP-код ошибки.