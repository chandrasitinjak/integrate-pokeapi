# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Install git (needed for go modules)
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o app ./main.go

# Stage 2: Run
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/app .
COPY .env .   
EXPOSE 8080
CMD ["./app"]
