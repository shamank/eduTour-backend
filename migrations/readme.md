### Создание миграции
```bash
migrate create -ext sql -dir ./migrations -seq *migration_title*
```

### Запуск миграций:
```bash
migrate -path ./migrations -database 'postgres://{user}:{password}@{host}:{port}/{db}?sslmode=disable' up
```

### Откат миграций:
```bash
migrate -path ./migrations -database 'postgres://{user}:{password}@{host}:{port}/{db}?sslmode=disable' down
```