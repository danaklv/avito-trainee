APP_NAME=pr-service
APP_CMD=./cmd/main.go

build:
	go build -o $(APP_NAME) $(APP_CMD)

run:
	docker-compose up --build

down:
	docker-compose down

lint:
	golangci-lint run ./...

test:
	go test ./... -v

load-test:
	k6 run k6-script.js
