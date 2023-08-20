up:
	docker-compose up --build app

test:
	go test --short -coverprofile=tests/coverage.out ./...
	make test.coverage

test.integration:
	docker run --rm -d -p 27018:27017 --name test -e MONGODB_DATABASE=test mongo:latest
	go test ./tests/
	docker stop test

test.coverage:
	go tool cover -html=tests/coverage.out

swag:
	swag init cmd/main.go

mock:
	mockgen -source=internal/service/service.go -destination=internal/service/mock/mock_service.go
	mockgen -source=internal/repository/repository.go -destination=internal/repository/mock/mock_repo.go