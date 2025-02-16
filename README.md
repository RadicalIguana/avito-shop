# Магазин мерча для сотрудников Авито

Этот проект представляет собой сервис для внутреннего магазина мерча компании Авито. Сотрудники могут использовать монеты для покупки товаров и передачи монет другим сотрудникам.

## Описание

Сервис позволяет:
- Покупать мерч за монеты.
- Передавать монеты другим сотрудникам.
- Просматривать историю транзакций и список купленных товаров.

## Стек технологий
- Язык программирования: Go
- База данных: PostgreSQL
- Авторизация: JWT
- Контейнеризация: Docker Compose
- Тестирование: Unit-тесты, интеграционные тесты

## Установка и запуск

### 1. Клонируйте репозиторий
```bash
git clone https://github.com/RadicalIguana/avito_shop.git
cd avito_shop
```

### 2. Соберите и запустите контейнеры
``` bash
docker-compose up --build
```
После выполнения команды сервис будет доступен по адресу: 
`http://localhost:8080`

### 3. Запуск миграций
База данных автоматически инициализируется при первом запуске контейнера. Миграции выполняются с использованием `golang-migrate`.
```bash
docker exec avito-shop-app-1 migrate "postgres:<DB_USER>:<DB_PASSWORD>@<DB_HOST>:<DB_PORT>/<DB_NAME>?sslmode=disable"
```

### 4. API
API сервиса соответствует спецификации, доступной в `schemes.yaml`.

### 5. Тестирование
Для запуска тестов выполните:
```bash
docker exec avito-shop-app-1 go test ./.. -cover
```
Тестовое покрытие превышает 40%. Включены юнит-тесты и интеграционные тесты.

### 6. Дополнительные задания
- Конфигурация линтера для Go находится в файле `.golangci.yaml`









### Нагрузочное тестирование
`vegeta attack -duration=10s -rate=100 -targets=targets.txt | vegeta report` 


Production db: `postgres://postgres:postgres@db:5432/avito_shop_db?sslmode=disable`
Test db в докере: `postgres://postgres:postgres@test_db:5433/test_db?sslmode=disable`
Test db локально: `postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable`






# Юнит-тесты
## 1. Запустить тесты и сгенерировать отчет
go test ./... -coverprofile=coverage.out

## 2. Вывести покрытие для каждой функции
go tool cover -func=coverage.out

## 3. Извлечь общий процент покрытия
go tool cover -func=coverage.out | grep total | awk '{print $3}'

## 4. Открыть HTML-отчет для визуализации
go tool cover -html=coverage.out