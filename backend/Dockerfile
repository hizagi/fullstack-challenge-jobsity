# Stage 1: Build Phase
FROM golang:1.23.2-alpine:3.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111Module=on GOAMD64=v4 go build -mod=vendor -o tasklist ./cmd/server

# Stage 2: Run Phase
FROM alpine:3.20

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/tasklist .

# Expose port 8080
EXPOSE 8080

# Command to run the app
CMD ["./tasklist"]
