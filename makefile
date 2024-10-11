# Database commands
db_up:
	docker-compose up postgres

db_up_d:
	docker-compose up postgres -d

db_down:
	docker-compose down postgres

# API commands
build:
	go build -o api ./cmd/api

run: build
	./api

docker_build:
	docker build -t your-api-image .

docker_run:
	docker run -p 8000:8000 --env-file .env your-api-image

run_app:
	docker-compose up

test:
	go test -v -cover ./...

.PHONY: db_up db_up_d db_down build run docker_build docker_run run_app test


