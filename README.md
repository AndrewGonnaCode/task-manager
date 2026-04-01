# Go Todo API - Microservices Architecture

Учебный проект на Go с микросервисной архитектурой и event-driven подходом для управления задачами (To-Do приложения) с интеграцией Kafka и Telegram.

## Архитектура

```
┌─────────────────┐         ┌──────────────┐         ┌────────────────────┐
│   HTTP Client   │────────▶│ task-service │────────▶│    PostgreSQL      │
│  (REST API)     │         │  (Port 8080) │         │                    │
└─────────────────┘         └──────────────┘         └────────────────────┘
                                    │
                                    │ Publish events
                                    ▼
                            ┌──────────────┐
                            │    Kafka     │
                            │ (task-events)│
                            └──────────────┘
                                    │
                                    │ Consume events
                                    ▼
                         ┌─────────────────────┐       ┌──────────────┐
                         │ notification-service│──────▶│  Telegram    │
                         │  (Kafka Consumer)   │       │     Bot      │
                         └─────────────────────┘       └──────────────┘
```

## Микросервисы

### 1. task-service
- **REST API** для CRUD операций с задачами
- **Kafka Producer** - публикует события при создании/обновлении/удалении задач
- **PostgreSQL** - хранение данных
- **Порт**: 8080

### 2. notification-service
- **Kafka Consumer** - слушает события из топика `task-events`
- **Telegram Bot** - отправляет уведомления о статусе задач
- Поддержка горизонтального масштабирования через consumer groups

## Функциональность

### Task Service API

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/health` | Health check |
| GET | `/api/v1/tasks` | Получить все задачи |
| GET | `/api/v1/tasks/:id` | Получить задачу по ID |
| POST | `/api/v1/tasks` | Создать новую задачу |
| PUT | `/api/v1/tasks/:id` | Обновить задачу |
| DELETE | `/api/v1/tasks/:id` | Удалить задачу |

### Kafka Events

- **task.created** - Задача создана
- **task.updated** - Задача обновлена
- **task.deleted** - Задача удалена

### Telegram Notifications

Notification-service отправляет в Telegram уведомления:
- ✅ Новая задача создана
- 📝 Задача обновлена
- ❌ Задача удалена

## Предварительные требования

- **Docker** и **Docker Compose**
- **Go 1.24.2+** (для локальной разработки)
- **Telegram Bot Token** (получить у @BotFather)
- **Make** (опционально, для удобных команд)

## Быстрый старт

### 1. Клонирование репозитория

```bash
git clone <repository-url>
cd go-todo-api
```

### 2. Настройка Telegram Bot

#### Получение Bot Token

1. Откройте Telegram и найдите **@BotFather**
2. Отправьте команду `/newbot`
3. Следуйте инструкциям для создания бота
4. Скопируйте полученный **bot token**

#### Получение Chat ID

1. Откройте Telegram и найдите **@userinfobot**
2. Отправьте любое сообщение
3. Скопируйте показанный **ID** (это ваш chat ID)

### 3. Настройка окружения

```bash
# Создать .env файл
cp .env.example .env

# Отредактировать .env и добавить ваши Telegram credentials
nano .env
```

Обновите следующие переменные в `.env`:
```bash
TELEGRAM_BOT_TOKEN=your_bot_token_here
TELEGRAM_CHAT_ID=your_chat_id_here
```

### 4. Запуск всех сервисов

```bash
# Использование Make (рекомендуется)
make up

# Или напрямую docker-compose
docker-compose up -d
```

### 5. Проверка статуса

```bash
# Проверить статус всех сервисов
make status

# Посмотреть логи
make logs
```

### 6. Тестирование API

```bash
# Health check
curl http://localhost:8080/health

# Создать задачу
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Первая задача",
    "description": "Тестирование микросервисной архитектуры"
  }'

# Получить все задачи
curl http://localhost:8080/api/v1/tasks

# Обновить задачу (замените {id} на ID задачи)
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "completed": true
  }'

# Удалить задачу
curl -X DELETE http://localhost:8080/api/v1/tasks/1
```

После каждой операции вы должны получить уведомление в Telegram!

## Команды Make

```bash
make help       # Показать все доступные команды
make build      # Собрать все сервисы локально
make test       # Запустить тесты для всех сервисов
make lint       # Запустить линтеры
make up         # Запустить все сервисы с docker-compose
make down       # Остановить все сервисы
make logs       # Показать логи всех сервисов
make restart    # Перезапустить все сервисы
make status     # Показать статус сервисов
make clean      # Очистить артефакты сборки и volumes
```

## Локальная разработка (без Docker)

### Запуск инфраструктуры

```bash
# Запустить только PostgreSQL, Kafka и Zookeeper
docker-compose up -d postgres kafka zookeeper kafka-init
```

### Запуск task-service

```bash
cd services/task-service
cp .env.example .env
# Отредактируйте .env при необходимости
go run cmd/main.go
```

### Запуск notification-service

```bash
cd services/notification-service
cp .env.example .env
# Добавьте Telegram credentials в .env
go run cmd/main.go
```

## CI/CD Pipeline

Проект использует **GitHub Actions** для автоматизации:

### Workflow

1. **Lint** - Проверка кода с golangci-lint
2. **Test** - Запуск тестов для всех сервисов
3. **Build** - Сборка Docker образов
4. **Deploy** - Деплой (симуляция для локального окружения)
5. **Telegram Notification** - Уведомление о статусе деплоя

### Триггеры

- **Push** в ветки `main`, `stage`, `dev` → полный pipeline
- **Pull Request** → lint и test

### Настройка GitHub Secrets

Добавьте следующие secrets в настройках репозитория:

```
DOCKER_HUB_USERNAME    - Ваш Docker Hub username
DOCKER_HUB_TOKEN       - Docker Hub access token
TELEGRAM_BOT_TOKEN     - Telegram bot token
TELEGRAM_CHAT_ID       - Telegram chat ID
```

**Путь**: Repository → Settings → Secrets and variables → Actions → New repository secret

## Структура проекта

```
go-todo-api/
├── services/
│   ├── task-service/              # REST API + Kafka Producer
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── models/            # Task & Event models
│   │   │   ├── handlers/          # HTTP handlers
│   │   │   ├── database/          # PostgreSQL connection
│   │   │   ├── kafka/             # Kafka producer
│   │   │   ├── middleware/        # CORS & Logger
│   │   │   └── config/            # Configuration
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   └── README.md
│   │
│   └── notification-service/      # Kafka Consumer + Telegram Bot
│       ├── cmd/main.go
│       ├── internal/
│       │   ├── models/            # Event models
│       │   ├── kafka/             # Consumer & Handler
│       │   ├── telegram/          # Telegram bot & Formatter
│       │   └── config/            # Configuration
│       ├── Dockerfile
│       ├── go.mod
│       └── README.md
│
├── .github/workflows/
│   └── deploy.yml                 # CI/CD pipeline
├── scripts/
│   └── setup.sh                   # Setup script
├── docker-compose.yml             # Infrastructure orchestration
├── Makefile                       # Common commands
├── .env.example                   # Environment template
└── README.md                      # This file
```

## Технологии

- **Go 1.24.2** - Язык программирования
- **Gin** - HTTP web framework
- **GORM** - ORM для PostgreSQL
- **Sarama** - Kafka client для Go
- **PostgreSQL 16** - Реляционная база данных
- **Kafka 7.6.0** - Message broker
- **Zookeeper** - Kafka coordination
- **Telegram Bot API** - Уведомления
- **Docker & Docker Compose** - Контейнеризация
- **GitHub Actions** - CI/CD

## Масштабирование

### Горизонтальное масштабирование

**task-service** (за load balancer):
```bash
docker-compose up -d --scale task-service=3
```

**notification-service** (consumer group):
```bash
docker-compose up -d --scale notification-service=3
```

### Вертикальное масштабирование

Увеличьте ресурсы контейнеров в `docker-compose.yml`:
```yaml
services:
  task-service:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
```

## Мониторинг и Логи

### Просмотр логов

```bash
# Все сервисы
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f task-service
docker-compose logs -f notification-service
docker-compose logs -f kafka
```

### Проверка Kafka топика

```bash
# Список топиков
docker-compose exec kafka kafka-topics --bootstrap-server localhost:9092 --list

# Описание топика task-events
docker-compose exec kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic task-events

# Чтение сообщений из топика
docker-compose exec kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic task-events --from-beginning
```

### Проверка PostgreSQL

```bash
# Подключение к БД
docker-compose exec postgres psql -U postgres -d todo_db

# SQL запросы
SELECT * FROM tasks;
```

## Troubleshooting

### Сервисы не запускаются

1. Проверьте, что Docker запущен
2. Проверьте доступность портов (8080, 5435, 9092)
3. Убедитесь, что `.env` файл существует и заполнен
4. Посмотрите логи: `docker-compose logs`

### Kafka ошибки

1. Подождите 30-60 секунд для инициализации Kafka
2. Проверьте health check Kafka: `docker-compose ps kafka`
3. Убедитесь, что Zookeeper запущен: `docker-compose ps zookeeper`

### Telegram уведомления не приходят

1. Проверьте `TELEGRAM_BOT_TOKEN` и `TELEGRAM_CHAT_ID` в `.env`
2. Убедитесь, что вы начали чат с ботом (отправьте `/start`)
3. Проверьте логи notification-service: `docker-compose logs notification-service`

### Task-service не подключается к БД

1. Проверьте, что PostgreSQL запущен: `docker-compose ps postgres`
2. Проверьте credentials в `.env`
3. Убедитесь, что порт 5435 не занят

## Production Considerations

Для production окружения рекомендуется:

1. **Managed Services**:
   - AWS RDS для PostgreSQL
   - AWS MSK или Confluent Cloud для Kafka
   - Отдельные Telegram боты для каждого окружения

2. **Security**:
   - SSL/TLS для всех соединений
   - Secrets management (AWS Secrets Manager, HashiCorp Vault)
   - Network isolation (VPC, security groups)

3. **Monitoring**:
   - Prometheus + Grafana для метрик
   - ELK stack для логов
   - Distributed tracing (Jaeger)

4. **Deployment**:
   - Kubernetes для оркестрации
   - Helm charts для управления
   - Blue-green или canary deployments

## Лицензия

Этот проект создан в образовательных целях.

## Автор

Создано с использованием микросервисной архитектуры и event-driven подхода.
