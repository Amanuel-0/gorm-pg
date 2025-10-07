FROM golang:1.25-alpine

# Install Air for live reload
RUN apk add --no-cache git \
    && go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod/go.sum & download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy rest of source
COPY . .

# Expose port
EXPOSE 8080
EXPOSE 5001

# Default command is set in docker-compose (air)