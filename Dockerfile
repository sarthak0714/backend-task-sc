# Use the official Golang image
FROM golang:1.22.2

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project
COPY . .

# Download dependencies
RUN go mod download

# Install Swagger packages
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go get -u github.com/swaggo/echo-swagger
RUN go get -u github.com/swaggo/files

# Generate Swagger documentation
RUN swag init -g cmd/main.go

# Build the application
RUN go build -o ./bin/server ./cmd/main.go

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./bin/server"]