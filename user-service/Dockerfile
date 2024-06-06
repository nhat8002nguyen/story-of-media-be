# Use the official golang image as the base image
FROM golang:1.21-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY ./src ./src

# Build the Go app
RUN go build -o main ./src/main.go

# Use scratch as the base image for a minimal container
FROM alpine:latest

# Set user permissions for security
RUN adduser -D myuser
USER myuser

# Copy the compiled Go binary from the build stage
COPY --from=build /app/main /app/main

# Expose the port the service will run on
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/app/main"]