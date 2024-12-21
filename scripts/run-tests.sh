# Dockerfile

# Use the official Golang image as a base
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Run tests
CMD ["go", "test", "./..."]