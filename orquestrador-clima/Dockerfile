# Estágio de Build
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Como o WORKDIR é /app e copiamos tudo, o build apontará para a pasta cmd/server relativa ao /app
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Estágio Final
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/server ./server

ENTRYPOINT ["./server"]