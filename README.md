# Видео 
Видео весит слишком много, чтобы загрузить прямо в markdown, поэтому залил на youtube
[https://youtu.be/K9bKpWStBJk](https://www.youtube.com/live/K9bKpWStBJk&feature=youtu.be)

## Установка и запуск сервиса

### Предварительные требования
- Установленный **Docker** и **Docker Compose**
- Установленный **Make** (опционально, для использования Makefile)

---

### Шаг 1: Создание `.env` файла
Создайте файл `.env` в корневой директории проекта с указанным содержимым.

```bash
cat > .env << EOF
SERVER_ADDR=0.0.0.0
SERVER_PORT=8081

POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_USER=user
POSTGRES_PASSWORD=12345
POSTGRES_DB=postgres_db

KAFKA_HEALTHCHECK_HOST=kafka
KAFKA_PORT=9092
KAFKA_HEALTHCHECK_TOPIC=__consumer_offsets
KAFKA_TOPIC=orders
EOF
```

### Шаг 2: Загрузка переменных окружения
Загрузите переменные окружения в текущую сессию.

```bash
source .env
```

### Шаг 3: Запуск сервиса
#### Вариант 1: Использование Makefile (рекомендуется)

```bash
# Запуск всех сервисов
make up

# Или пошаговый запуск:
# Сначала запустите инфраструктуру
make up-without-app

# Затем запустите приложение
make up-app
```

#### Вариант 2: Прямое использование Docker Compose
- Запуск всех сервисов.
- Или запуск с пересборкой приложения.

```bash
# Запуск всех сервисов
docker compose up -d

# Или с пересборкой приложения
docker compose up --build -d
```

### Доступ к сервисам
- Приложение: [http://localhost:8081](http://localhost:8081)
- Kafka UI: [http://localhost:9020](http://localhost:9020)
- PostgreSQL: `localhost:5432`

## Схема базы данных 


![img.png](img.png)
