FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@latest


COPY go.mod go.sum ./
RUN go mod download

COPY .air.toml ./

COPY . .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]

