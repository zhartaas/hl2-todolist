run:
	go run cmd/web/*

swagger:
	swag init -g cmd/web/main.go -o ./docs