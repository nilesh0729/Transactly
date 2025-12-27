# Build Stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main cmd/api/main.go

# Run Stage

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD [ "/app/main" ]
