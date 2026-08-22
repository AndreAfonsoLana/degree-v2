FROM golang:1.25.0 AS builder
WORKDIR /app
COPY . .
RUN go mod download

# 🎯 Ajuste exato baseado na sua estrutura de pastas (cmd/server)
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server/main.go

# Estágio 2: Execução
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]