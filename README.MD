
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