FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum./

RUN go mod init RE

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 8025

CMD ["./main"]



