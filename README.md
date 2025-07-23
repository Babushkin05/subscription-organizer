# 📦 Subscription Organizer - тестовое задание для Effective Mobile

Subscription Organizer — это сервис на Go, предназначенный для учета пользовательских подписок и расчета их стоимости за указанный период.

Проект следует **чистой архитектуре** с разделением на домен, приложения, инфраструктуру и общий слой.

---

## 📁 Структура проекта

```
.
├── cmd/                            # Точка входа (main.go)
├── internal/
│   ├── application/
│   │   ├── service/               # Бизнес-логика (Application Layer)
│   │   └── port/                  # Интерфейсы для service и repository
│   ├── config/                    # Загрузка и парсинг конфигурации
│   ├── domain/
│   │   └── model/                 # Доменные модели
│   ├── infrastructure/
│   │   ├── delivery/
│   │   │   └── http/              # HTTP хендлеры (Gin)
│   │   └── repository/           # Реализация репозиториев (PostgreSQL)
│   └── shared/
│       ├── dto/                   # DTO объекты запроса/ответа
│       └── mapper/                # Преобразование DTO <-> Model
├── migrations/                    # SQL миграции (golang-migrate)
├── pkg/
│   └── logger/                    # Кастомный логгер
├── docs/                          # Сгенерированная swagger-документация
├── .env                           # Переменные окружения
├── Makefile                       # Утилиты для сборки, миграций и т.д.
├── docker-compose.yaml            # Контейнеризация приложения и БД
└── go.mod / go.sum                # Зависимости Go
```

---

## ⚙️ Запуск проекта

### 1. Клонируйте репозиторий

```bash
git clone https://github.com/Babushkin05/subscription-organizer.git
cd subscription-organizer
```

### 2. Настройка `.env`

Создайте файл `.env` на основе шаблона:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
SERVER_PORT=8080
```

### 3. Сборка и запуск через Docker

```bash
make up
```

### 4. Выполнение миграций

```bash
make migrate-up
```

### 4. Генерация swagger-doc

```bash
make swaga
```

### 5. Запуск через docker-compose

```bash
docker-compose up --build 
```

После запуска:

```
http://localhost:8080/swagger/index.html
```

---

## 🧱 Структуры запросов

```json
// CreateSubscriptionRequest (пример запроса на создание подписки)
{
  "service_name": "Netflix",
  "price": 999,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "07-2025",
  "end_date": "12-2025"
}

// SubscriptionResponse (пример успешного ответа с подпиской)
{
  "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "service_name": "Netflix",
  "price": 999,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "07-2025",
  "end_date": "12-2025"
}

// ErrorResponse (пример ответа при ошибке)
{
  "error": "invalid user_id"
}

// MessageResponse (пример сообщения об успешном действии)
{
  "message": "subscription deleted"
}

// CalculateTotalCost Response (пример ответа при подсчёте общей стоимости)
{
  "total_cost": 5998
}

```


---

## 🧪 Примеры запросов (curl)

```bash
# Создание подписки
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Netflix",
    "price": 1000,
    "user_id": "user-uuid-here",
    "start_date": "07-2025",
    "end_date": "12-2025"
}'

# Получение всех подписок
curl http://localhost:8080/subscriptions

# Расчет общей стоимости
curl "http://localhost:8080/subscriptions/cost?user_id=user-uuid&from=01-2025&to=12-2025"
```

---

## 🧰 Makefile команды

```bash
make up              # Запустить контейнеры
make down            # Остановить контейнеры
make migrate-up      # Применить миграции
make migrate-down    # Откатить миграции
make swag            # Сгенерировать Swagger документацию
```

---

## 📌 Используемые технологии

- Go 1.21+
- Gin
- PostgreSQL
- Docker + docker-compose
- golang-migrate
- swaggo/swag (Swagger)
- Zap Logger


