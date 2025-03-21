# Estágio de construção
FROM golang:1.24-alpine AS builder

# Define o ambiente de trabalho
WORKDIR /app

# RUN apk add --no-cache git ca-certificates upx && \
#     rm -rf /var/cache/apk/*

# Instala dependências mínimas necessárias
# RUN apk add --no-cache upx && \
#     rm -rf /var/cache/apk/*

RUN apk add --no-cache upx && \
    rm -rf /var/cache/apk/*

# Copia apenas os módulos para aproveitar o cache
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copia o código-fonte
COPY . .

# Compila com otimizações máximas e comprime com UPX
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
#     -trimpath \
#     -ldflags="-s -w -extldflags '-static' -buildid=" \
#     -buildvcs=false \
#     -o /main ./cmd && \
#     upx --best /main

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -extldflags '-static' -buildid=" \
    -buildvcs=false \
    -o /main ./cmd

# Estágio final minimalista
FROM scratch
COPY --from=builder /main /main
EXPOSE 8080
ENTRYPOINT ["/main"]