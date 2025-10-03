FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /app/main .

EXPOSE 4000

CMD ["./main"]


FROM builder AS test

CMD ["go", "test", "-v", "./cmd/api/..."]
