FROM golang:alpine
WORKDIR /app
EXPOSE 44044
CMD go run cmd/sso/main.go
