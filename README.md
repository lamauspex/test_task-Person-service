# Person Service

> RESTful API для управления данными сотрудников на Go

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Echo](https://img.shields.io/badge/Echo-4.13.4-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)
![Swagger](https://img.shields.io/badge/Swagger-Yes-green)

## 🟢 Технологический стек

| Компонент | Технология | Назначение |
|-----------|------------|------------|
| Язык | Go 1.21+ | Основной язык разработки |
| Web-фреймворк | Echo v4 | HTTP-маршрутизация и обработчики |
| База данных | PostgreSQL | Хранение данных |
| Логирование | Uber Zap | Структурированное логирование |
| Валидация | go-playground/validator | Валидация входных данных |
| Документация | Swaggo/Swagger | Автогенерация API-документации |
| Конфигурация | Godotenv | Работа с .env файлами |
| JSON | ByteDance Sonic | Высокопроизводительная сериализация |

## 🟢 Ключевые навыки

- **RESTful API** — проектирование и реализация
- **Clean Architecture** — разделение на слои (handler → logic → repo)
- **PostgreSQL** — работа с реляционной БД
- **Логирование** — структурированные логи с Zap
- **Валидация данных** — проверка входящих запросов
- **Docker** — готовность к контейнеризации (опционально)
- **Swagger** — автоматическая документация API

## 🟢 Возможности

- Создание нового сотрудника (Person)
- Получение списка всех сотрудников
- Получение сотрудника по ID
- Обновление данных сотрудника
- Удаление сотрудника
- Валидация входящих данных
- Логирование всех запросов
- Swagger-документация

## 🟢 Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/persons` | Получить всех сотрудников |
| `GET` | `/persons/:id` | Получить сотрудника по ID |
| `POST` | `/persons` | Создать нового сотрудника |
| `PUT` | `/persons/:id` | Обновить сотрудника |
| `DELETE` | `/persons/:id` | Удалить сотрудника |

## 🟢 Примеры запросов

### Создание сотрудника

```bash
curl -X POST http://localhost:8080/persons \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "phone": "+79001234567",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Ответ

```json
{
  "id": 1,
  "email": "john.doe@example.com",
  "phone": "+79001234567",
  "first_name": "John",
  "last_name": "Doe"
}
```

### Получение списка

```bash
curl http://localhost:8080/persons
```

## 🟢 Структура проекта

```
Person_service/
├── main.go                    # Точка входа
├── .env                       # Конфигурация
├── internal/
│   ├── app/                   # Модели данных
│   │   └── person_model.go
│   ├── http/                  # HTTP-обработчики
│   │   └── person_handler.go
│   ├── logic/                 # Бизнес-логика
│   │   └── person_logic.go
│   ├── db/                    # Работа с БД
│   │   └── postgre_person_repo.go
│   ├── middleware/            # Middleware
│   │   └── logger_middleware.go
│   └── tests/                 # Интеграционные тесты
├── migrations/                # SQL-миграции
├── docs/                      # Swagger-документация
└── go.mod / go.sum            # Зависимости
```

## 🟢 Быстрый старт

```bash
# Клонирование
git clone https://github.com/lamauspex/test_task-Person-service

# Переход в директорию
cd Person_service

# Настройка .env
# Отредактируйте DB_PASSWORD в файле .env

# Установка зависимостей
go mod download

# Запуск
go run main.go
```

## 🟢 Документация

После запуска доступна по адресу:

```
http://localhost:8080/swagger/index.html
```

## 🟢 Архитектура

```
┌─────────────┐     ┌─────────────┐     ┌──────────────┐
│   Handler   │ ──► │   Service   │ ──► │ Repository   │
│   (Echo)    │     │   (Logic)   │     │  (PostgreSQL)│
└─────────────┘     └─────────────┘     └──────────────┘
       │                   │                   │
       ▼                   ▼                   ▼
  HTTP Request      Business Logic       SQL Queries
```

- **Handler** — принимает HTTP-запросы, валидирует, передаёт в сервис
- **Service** — бизнес-логика, обработка данных
- **Repository** — доступ к базе данных, SQL-запросы

---

**Автор**: Резник Кирилл  
**Email**: lamauspex@yandex.ru  
**Telegram**: @lamauspex  
**GitHub**: https://github.com/lamauspex
