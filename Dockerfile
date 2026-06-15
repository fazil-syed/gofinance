FROM golang:1.25-alpine AS builder

WORKDIR /temp

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN go build -o app ./cmd/main/gofinance.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /temp/app .

COPY config.yaml .

EXPOSE 8080

CMD [ "./app" ]