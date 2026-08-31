# Multi-stage build for Go Backend (Atlas Food)
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/api/main.go

# Production minimal runtime image
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server ./server
COPY --from=builder /app/Atlas_Makananku_FINAL.json ./Atlas_Makananku_FINAL.json

RUN mkdir -p uploads

EXPOSE 8080

CMD ["./server"]
