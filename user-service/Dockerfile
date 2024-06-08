# Use the official golang image as the base image
FROM golang:1.21-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY .env .env

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY src/ src/

# Build the Go app
RUN go build -o main ./src

# Use scratch as the base image for a minimal container
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary and .env from the builder stage
COPY --from=build /app/main .
COPY --from=build /app/.env .env

# Set permissions for the .env file (optional but good practice)
RUN chmod 644 .env

# Indicate the port that the container app will listen on
EXPOSE 8081

# Command to run the executable
ENTRYPOINT ["./main"]