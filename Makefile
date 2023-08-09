run:
	go run cmd/main.go

up:
	docker-compose up --build app

test:
	go test ./...