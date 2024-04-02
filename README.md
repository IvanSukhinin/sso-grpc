# SSO GRPC API

## local deploy
1) Containers build and up
```shell
docker-compose up -d
```
2) Migrations up
```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
```

```shell
cd ./migrations && \
goose postgres "host=localhost port=54322 user=postgres password=postgres database=sso sslmode=disable" up
```
