# Trackly Backend

Бэкенд-часть проекта [Trackly](https://github.com/All4A/Trackly), реализованная на Go.

## Техническая реализация (Backend)

- **Go** — с использованием `Echo` (API) и `Gorm` (ORM)
- **oapi-codegen** — генерация эндпоинтов по OpenAPI
- **golang-migrate** — миграции схемы БД
- **golang-jwt** — middleware для JWT-аутентификации
- **gopkg.in/yaml.v3** — конфигурация из YAML
- **PostgreSQL** + **Minio S3** — хранилище данных и файлов

## 🛠 Установка и запуск

### 1. Очистка кеша и установка зависимостей

```bash
go clean -modcache
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest


## Codegen
используется кодогенерация ручек по опен апи схеме [openapi.yaml](open-api/openapi.yaml)

Чистим кеш 
```bash
go clean -modcache
```

Устанавливаем зависимости
```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```
Генерим коды
```bash
oapi-codegen -generate client,types,server -package api open-api/openapi.yaml > internal/api/api.gen.go
```


Запускается все дело в докере, вот эту штуку можно запустить сразу
```bash
docker-compose up --build
```
