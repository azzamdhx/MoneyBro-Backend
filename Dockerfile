FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrate

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /server /app/server
COPY --from=builder /migrate /app/migrate
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/email-templates /app/email-templates

ENV TZ=Asia/Jakarta

EXPOSE 8080

CMD ["/bin/sh", "-c", "/app/migrate up && /app/server"]
