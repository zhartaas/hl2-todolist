up:
	docker-compose up

up-build:
	docker-compose up --build

up-deattached:
	docker-compose up -d

up-build-deattached:
	docker-compose up --build -d

build:
	docker-compose build

down:
	docker-compose down

run-local:
	go run cmd/web/*

swagger:
	swag init -g cmd/web/main.go -o ./docs