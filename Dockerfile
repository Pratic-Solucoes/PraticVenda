# Stage 1: Build the Go binary
# Usamos uma imagem oficial do Go com Alpine Linux, que é leve.
FROM golang:1.25-alpine AS builder

# Instala git, necessário para o 'go mod download' buscar dependências de repositórios.
RUN apk add --no-cache git

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos de gerenciamento de dependências primeiro para aproveitar o cache do Docker.
# Se esses arquivos não mudarem, o Docker não baixará as dependências novamente.
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o código-fonte do projeto para o contêiner.
COPY . .

# Compila a aplicação Go.
# CGO_ENABLED=0: Cria um binário estático, sem dependências de bibliotecas C do sistema.
# -ldflags '-w -s': Remove informações de debug, resultando em um binário menor.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s' -o gestao-api ./cmd/api

# Stage 2: Create the final, minimal production image
# Usamos a imagem base do Alpine, que é uma das menores disponíveis (~5MB).
FROM alpine:latest

# Instala certificados CA para que sua aplicação Go possa fazer chamadas HTTPS seguras.
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copia o binário compilado e a pasta 'web' (com templates e assets) do estágio 'builder'.
COPY --from=builder /app/gestao-api .
COPY --from=builder /app/web ./web

# Expõe a porta em que a aplicação será executada.
EXPOSE 9000

# Comando para iniciar a aplicação quando o contêiner for executado.
CMD ["./gestao-api"]