# Use a Go base image for development
FROM golang:1.24-alpine

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Install air
RUN go install github.com/air-verse/air@latest

# Copy your .air.toml configuration
# This should be in the same directory as your Dockerfile (i.e., backend/.air.toml)
COPY .air.toml .
COPY ../.env .

# Expose the port the application listens on
EXPOSE 8080

# The CMD here sets the default command if not overridden.
# In your docker-compose.yml, you *are* overriding it with 'command: sh -c "sleep 5 && air"'
# So this CMD is less critical here but good practice.
CMD ["air", "-c", ".air.toml"]
