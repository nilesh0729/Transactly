# Build Stage
FROM golang:1.25.2-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .

EXPOSE 8080
CMD [ "/app/main" ]
