# Trackly Backend

Бэкенд-часть проекта [Trackly](https://github.com/All4A/Trackly), реализованная на Go.

## Техническая реализация (Backend)

- **Go** — с использованием `Echo` (API) и `Gorm` (ORM)
- **oapi-codegen** — генерация эндпоинтов по OpenAPI
- **golang-migrate** — миграции схемы БД
- **golang-jwt** — middleware для JWT-аутентификации
- **gopkg.in/yaml.v3** — конфигурация из YAML
- **PostgreSQL** + **Minio S3** — хранилище данных и файлов
