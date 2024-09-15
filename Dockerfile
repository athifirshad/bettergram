# Build stage
FROM golang:1.23-alpine AS build

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api
# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Create uploads directory
RUN mkdir -p /uploads && chmod 777 /uploads

# Copy the compiled application from the build stage
COPY --from=build /app/main .

# Expose the port the app runs on
EXPOSE 4000

# Run the application
CMD ["./main"]