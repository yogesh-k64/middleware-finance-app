# Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o finance-app

EXPOSE 9000

CMD ["./finance-app"]