# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /app 

COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 5000
CMD ["/app/main"]