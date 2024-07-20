FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ./hl2-todolist ./cmd/web

CMD ["./hl2-todolist"]