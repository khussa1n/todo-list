up:
	docker-compose up --build app

test:
	go test --short -coverprofile=tests/coverage.out ./...
	make test.coverage

test.integration.with.dockerDB:
	docker run --rm -d -p 27018:27017 --name test -e MONGODB_DATABASE=test mongo:latest
	go test ./tests/
	dockr stop test
	#dockr rm test

test.integration.with.existing.dockerDB:
	go test  ./tests/

test.coverage:
	go tool cover -html=tests/coverage.out

swag:
	swag init cmd/main.go