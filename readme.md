if you are running without docker, you can use

```
go install github.com/joho/godotenv/cmd/godotenv@latest
```

```
godotenv go run .
```

where godotenv loads `.env` file.

**migrate database**

```bash
migrate -path migrations -database postgres://postgres:postgres@postgres/sql_rest_test0001?sslmode=disable up
```

creating volume `postgres_test_go` so data will be saved.

```bash
docker volume create postgres_test_go
```

and you can use `external: true` in docker-compose.yaml so if you use --down -v it wont remove data

**Launch app**

```bash
docker-compose up --build
```

**OR** (if you want to run test from this)

```bash
docker-compose build
docker-compose up -d postgres
```


**Tests**

``godotenv go test ./routes``
