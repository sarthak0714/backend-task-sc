# Use the official Golang image
FROM golang:1.22.2

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o ./bin/server ./cmd/main.go

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./bin/server"]