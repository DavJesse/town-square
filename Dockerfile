# Set up the latest Golang base image
FROM golang:1.23.0

# Set the working directory inside the container
WORKDIR /app

# Copy go module files first for efficient caching
COPY go.mod go.sum ./

# Download Go dependencies
RUN go mod download

# Copy the rest of the project files to the current working directory
COPY . .

# Set environment variables
ENV APP_NAME="forum"
ENV APP_PORT=8080
ENV DB_HOST="localhost"
ENV DB_NAME="forum.db"

# Build the application binary
RUN go build -o forum

# Expose the application port
EXPOSE 8080

# Default command to run the application
CMD ["./forum"]
