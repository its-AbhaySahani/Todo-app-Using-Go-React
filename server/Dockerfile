# Use the official Golang image as the base image
FROM golang:1.20 as build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Use a minimal base image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Install necessary packages
RUN apk --no-cache add ca-certificates

# Copy the built Go application from the build stage
COPY --from=build /app/main .

# Copy the MySQL queries
COPY --from=build /app/Database/MySql_Queries/query.sql ./Database/MySql_Queries/query.sql

# Expose the application port
EXPOSE 9000

# Command to run the application
CMD ["./main"]