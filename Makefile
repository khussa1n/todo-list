build:
	docker-compose build app

run:
	docker-compose up --build app

test:
	go test ./...