.PHONY: proto build run test clean

# Генерация протобуф файлов
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/notification.proto

# Сборка проекта
build:
	go build -o bin/notifyhub cmd/notifyhub/main.go

# Запуск сервера
run:
	go run cmd/notifyhub/main.go

# Тесты
test:
	go test ./...

# Очистка
clean:
	rm -rf bin/

# Установка зависимостей
deps:
	go mod tidy
	go mod download

# Docker сборка
docker-build:
	docker build -t notify-hub .