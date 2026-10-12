# Stage 1: Build binary
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
# Download dependency dan build binary
RUN go mod download
RUN go build -o otp-engine cmd/api/main.go

# Stage 2: Minimalist image untuk production
FROM alpine:latest
WORKDIR /app
# Copy hasil build dari stage 1
COPY --from=builder /app/otp-engine .
# Copy file .env jika ada
COPY .env .env 

EXPOSE 4000
CMD ["./otp-engine"]
