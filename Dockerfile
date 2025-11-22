FROM golang:1.23.1 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o service ./cmd/


FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/service .
COPY migrations ./migrations
COPY entrypoint.sh .

RUN apt-get update && apt-get install -y postgresql-client
RUN chmod +x entrypoint.sh

EXPOSE 8080

CMD ["./entrypoint.sh"]
