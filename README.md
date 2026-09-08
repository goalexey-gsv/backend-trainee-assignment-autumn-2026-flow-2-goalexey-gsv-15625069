# Kitchen API

MVP backend для сервиса доставки еды из ресторанов.

Проект предоставляет REST API для клиентской части и ресторанов, хранит данные в PostgreSQL и взаимодействует с отдельным сервисом ресторана.

## Возможности

### Пользователь

- просмотр списка доступных ресторанов;
- просмотр меню ресторана;
- создание заказа;
- просмотр заказа и его текущего статуса;
- просмотр истории заказов пользователя.

### Ресторан

- получение заказа от основного API;
- подтверждение или отклонение заказа;
- изменение статуса заказа в рамках жизненного цикла.

## Архитектура

Проект состоит из двух Go-сервисов:

```text
                         ┌──────────────────────┐
                         │       Client         │
                         │   Web / Mobile UI    │
                         └──────────┬───────────┘
                                    │
                                    │ REST API
                                    ▼
                    ┌──────────────────────────────┐
                    │          Kitchen API         │
                    │                              │
                    │ Go + Gin                    │
                    │                              │
                    │ ┌────────┐ ┌──────────────┐ │
                    │ │Handler │→│   Service    │ │
                    │ └────────┘ └──────┬───────┘ │
                    │                   │         │
                    │              ┌────▼─────┐   │
                    │              │Repository│   │
                    │              └────┬─────┘   │
                    └───────────────────┼─────────┘
                                        │
                              ┌─────────▼─────────┐
                              │    PostgreSQL     │
                              └───────────────────┘

                    HTTP REST              HTTP REST
                         │                     │
                         ▼                     ▼
              ┌─────────────────┐    ┌────────────────────┐
              │ Kitchen API     │───→│ Restaurant Mock    │
              │                 │    │                    │
              │ Order Service   │    │ Go + Gin           │
              └─────────────────┘    └────────────────────┘
```

### Основной сервис

`Kitchen API` отвечает за:

- REST API;
- бизнес-логику;
- работу с PostgreSQL;
- создание и хранение заказов;
- проверку состава заказа;
- взаимодействие с сервисом ресторана.

### Restaurant Mock

`restaurant-mock` — отдельный сервис, имитирующий API реального ресторана.

Он запускается вместе с основным приложением через Docker Compose и используется для демонстрации межсервисного взаимодействия.

## Архитектурный подход

В основном сервисе используется разделение на:

```text
Handler
   ↓
Service
   ↓
Repository
   ↓
PostgreSQL
```

Для взаимодействия с внешним сервисом ресторана используется отдельный HTTP client:

```text
OrderService
     ↓
Restaurant Client
     ↓ HTTP
Restaurant Mock
```

Это позволяет отделить HTTP-слой, бизнес-логику, доступ к данным и интеграции.

## Структура проекта

```text
.
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── handler/
│   │   ├── order.go
│   │   └── restaurant.go
│   │
│   ├── integration/
│   │   └── restaurant/
│   │       └── client.go
│   │
│   ├── model/
│   │   ├── menu_item.go
│   │   ├── order.go
│   │   └── restaurant.go
│   │
│   ├── repository/
│   │   ├── order.go
│   │   ├── restaurant.go
│   │   └── postgres/
│   │       ├── order.go
│   │       └── restaurant.go
│   │
│   └── service/
│       ├── order.go
│       └── restaurant.go
│
├── migrations/
│   ├── 00001_init.sql
│   ├── 00002_orders.sql
│   └── 00003_seed.sql
│
├── docs/
│   ├── cjm/
│   │   ├── user-cjm.puml
│   │   ├── user-cjm.png
│   │   ├── restaurant-cjm.puml
│   │   └── restaurant-cjm.png
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── restaurant-mock/
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
│
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── .golangci.yml
├── go.mod
└── go.sum
```

## Запуск

### Требования

- Go 1.26+
- Docker
- Docker Compose

Для локальной разработки также используются:

- `swag`;
- `plantuml`;
- `golangci-lint`.

### Запуск через Docker Compose

Основной способ запуска проекта:

```bash
docker compose up --build
```

После запуска:

```text
Kitchen API:       http://localhost:8080
Restaurant Mock:   http://localhost:8081
PostgreSQL:        localhost:5432
```

Health check:

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```json
{
  "status": "ok"
}
```

Остановить сервисы:

```bash
docker compose down
```

Для удаления данных PostgreSQL вместе с volume:

```bash
docker compose down -v
```

Миграции базы данных запускаются автоматически при старте `Kitchen API`.

## API

API использует версионирование:

```text
/api/v1
```

### Restaurants

Получить список активных ресторанов:

```http
GET /api/v1/restaurants
```

Получить ресторан:

```http
GET /api/v1/restaurants/{id}
```

Получить доступное меню:

```http
GET /api/v1/restaurants/{id}/menu
```

### Orders

Создать заказ:

```http
POST /api/v1/orders
```

Пример:

```json
{
  "restaurant_id": "11111111-1111-1111-1111-111111111111",
  "delivery_address": "Москва, ул. Ленина, 15",
  "items": [
    {
      "menu_item_id": "21111111-1111-1111-1111-111111111111",
      "quantity": 2
    },
    {
      "menu_item_id": "23333333-3333-3333-3333-333333333333",
      "quantity": 1
    }
  ]
}
```

Идентификатор пользователя передаётся через заголовок:

```http
X-User-ID: 99999999-9999-9999-9999-999999999999
```

Получить заказ:

```http
GET /api/v1/orders/{id}
```

Получить заказы пользователя:

```http
GET /api/v1/users/{userId}/orders
```

## Жизненный цикл заказа

```text
pending
   │
   ├──────────────→ rejected
   │
   ▼
accepted
   │
   ▼
preparing
   │
   ▼
ready
   │
   ▼
delivering
   │
   ▼
delivered
```

Также предусмотрен статус:

```text
cancelled
```

При создании заказа основной API:

1. проверяет ресторан и позиции меню;
2. проверяет наличие товаров;
3. фиксирует актуальные цены в `order_items`;
4. рассчитывает итоговую стоимость;
5. создаёт заказ со статусом `pending`;
6. отправляет заказ в `restaurant-mock`;
7. получает результат от ресторана;
8. обновляет статус заказа.

Цена в базе хранится в копейках в типе `BIGINT`.

Например:

```text
45000 = 450.00 ₽
```

В `order_items` сохраняются `name` и `price` на момент создания заказа. Это позволяет сохранить историческое состояние заказа даже после изменения меню.

## Database Schema

Основные сущности:

```text
restaurants
     │
     │ 1:N
     ▼
menu_items

restaurants
     │
     │ 1:N
     ▼
orders
     │
     │ 1:N
     ▼
order_items
```

### restaurants

| Поле        | Тип         | Описание                                   |
| ----------- | ----------- | ------------------------------------------ |
| id          | UUID        | Идентификатор                              |
| external_id | VARCHAR     | Идентификатор ресторана во внешней системе |
| name        | VARCHAR     | Название                                   |
| description | TEXT        | Описание                                   |
| address     | TEXT        | Адрес                                      |
| status      | VARCHAR     | `active` / `inactive`                      |
| created_at  | TIMESTAMPTZ | Дата создания                              |
| updated_at  | TIMESTAMPTZ | Дата обновления                            |

### menu_items

| Поле          | Тип         | Описание                                 |
| ------------- | ----------- | ---------------------------------------- |
| id            | UUID        | Идентификатор                            |
| restaurant_id | UUID        | Ресторан                                 |
| external_id   | VARCHAR     | Идентификатор позиции во внешней системе |
| name          | VARCHAR     | Название                                 |
| description   | TEXT        | Описание                                 |
| price         | BIGINT      | Цена в копейках                          |
| available     | BOOLEAN     | Доступность                              |
| created_at    | TIMESTAMPTZ | Дата создания                            |
| updated_at    | TIMESTAMPTZ | Дата обновления                          |

### orders

| Поле             | Тип         | Описание                      |
| ---------------- | ----------- | ----------------------------- |
| id               | UUID        | Идентификатор заказа          |
| restaurant_id    | UUID        | Ресторан                      |
| user_id          | UUID        | Идентификатор пользователя    |
| status           | VARCHAR     | Статус заказа                 |
| total_price      | BIGINT      | Итоговая стоимость в копейках |
| delivery_address | TEXT        | Адрес доставки                |
| created_at       | TIMESTAMPTZ | Дата создания                 |
| updated_at       | TIMESTAMPTZ | Дата обновления               |

### order_items

| Поле         | Тип     | Описание                  |
| ------------ | ------- | ------------------------- |
| id           | UUID    | Идентификатор             |
| order_id     | UUID    | Заказ                     |
| menu_item_id | UUID    | Позиция меню              |
| name         | VARCHAR | Название на момент заказа |
| price        | BIGINT  | Цена на момент заказа     |
| quantity     | INTEGER | Количество                |

## CJM

В проекте представлены code-generated CJM-диаграммы в формате PlantUML.

### Пользователь

![User CJM](docs/cjm/user-cjm.png)

Исходник:

```text
docs/cjm/user-cjm.puml
```

### Ресторан

![Restaurant CJM](docs/cjm/restaurant-cjm.png)

Исходник:

```text
docs/cjm/restaurant-cjm.puml
```

Диаграммы описывают основные пользовательские сценарии:

```text
Пользователь:
ресторан → меню → корзина → заказ → принятие/отклонение → доставка

Ресторан:
получение заказа → проверка → принятие/отклонение → приготовление → доставка
```

## OpenAPI / Swagger

Swagger UI доступен после запуска приложения:

```text
http://localhost:8080/swagger/index.html
```

Исходная спецификация генерируется из Go-аннотаций.

Для генерации:

```bash
make swagger
```

Команда использует:

```bash
swag init -g cmd/api/main.go -o docs
```

В результате обновляются:

```text
docs/docs.go
docs/swagger.json
docs/swagger.yaml
```

## Development

### Makefile

Доступные команды:

```bash
make swagger
```

Генерация Swagger-документации.

```bash
make diagrams
```

Генерация CJM-диаграмм из PlantUML.

```bash
make lint
```

Запуск `golangci-lint`.

```bash
make build
```

Сборка проекта.

```bash
make test
```

Запуск тестов.

```bash
make run
```

Запуск API локально.

```bash
make docker-up
```

Сборка и запуск всех сервисов через Docker Compose.

```bash
make docker-down
```

Остановка Docker Compose.

## Проверка сценария

После запуска проекта можно проверить основной сценарий.

### 1. Получить рестораны

```bash
curl http://localhost:8080/api/v1/restaurants
```

### 2. Получить меню

```bash
curl http://localhost:8080/api/v1/restaurants/11111111-1111-1111-1111-111111111111/menu
```

### 3. Создать заказ

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 99999999-9999-9999-9999-999999999999" \
  -d '{
    "restaurant_id": "11111111-1111-1111-1111-111111111111",
    "delivery_address": "Москва, ул. Ленина, 15",
    "items": [
      {
        "menu_item_id": "21111111-1111-1111-1111-111111111111",
        "quantity": 2
      }
    ]
  }'
```

Ресторан `restaurant-mock` принимает заказ и возвращает:

```json
{
  "external_order_id": "restaurant-<order_id>",
  "status": "accepted"
}
```

После этого основной API переводит заказ из:

```text
pending → accepted
```

### 4. Получить заказ

```bash
curl http://localhost:8080/api/v1/orders/{order_id}
```

## Технологии

- Go 1.26
- Gin
- PostgreSQL 17
- pgx
- Goose
- Swagger / OpenAPI
- PlantUML
- Docker / Docker Compose
- golangci-lint

## Ограничения MVP

Проект намеренно ограничен рамками тестового MVP.

Не реализованы:

- аутентификация и авторизация;
- реальная платёжная система;
- реальная служба доставки;
- полноценная админ-панель;
- push/email-уведомления;
- полноценная синхронизация меню с внешними ресторанами;
- распределённая обработка событий.

`restaurant-mock` используется вместо реальной внешней системы ресторана, чтобы продемонстрировать интеграцию сервисов в рамках одного репозитория.

## Возможное развитие: разграничение доступа между сервисами

В текущем MVP аутентификация и авторизация намеренно не реализованы, так как они не входят в обязательные требования задания.

В production-версии взаимодействие между сервисами можно дополнительно защитить с помощью JWT.

Предлагаемый вариант:

1. Сервис получает JWT в заголовке `Authorization: Bearer <token>`.
2. Middleware основного API выполняет валидацию подписи и срока действия токена.
3. Из claims токена извлекается идентификатор сервиса (`service_id`).
4. `service_id` сверяется со списком сервисов, которым разрешён доступ к конкретному API.
5. Для каждого endpoint можно определить требуемые права или список допустимых сервисов.
6. При отсутствии или некорректности токена запрос отклоняется с `401 Unauthorized`, а при отсутствии необходимых прав — с `403 Forbidden`.

Например:

```text
JWT
 ├── service_id: kitchen-api
 ├── permissions:
 │    ├── orders.read
 │    └── orders.create
 └── exp: ...

             │
             ▼
      JWT Middleware
             │
             ▼
     Validate signature
             │
             ▼
       Extract service_id
             │
             ▼
   Check allowed services
             │
       ┌─────┴─────┐
       │           │
     allow        deny
       │           │
       ▼           ▼
   Handler      403 Forbidden
```

Список разрешённых сервисов можно хранить в конфигурации либо в отдельной таблице БД:

```text
service_permissions

service_id       resource          action
------------------------------------------------
restaurant-001   orders            read
restaurant-001   orders            update_status
kitchen-api      orders            create
kitchen-api      restaurants       read
```

Такой подход позволит в дальнейшем перейти от демонстрационного `restaurant-mock` к нескольким независимым сервисам с изолированными правами доступа, не меняя основную бизнес-логику приложения.
