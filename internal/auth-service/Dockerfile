FROM golang:alpine

WORKDIR /app

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 8025

CMD ["./main"]



