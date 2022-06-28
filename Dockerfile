# syntax=docker/dockerfile:1
FROM golang:1.17-alpine

# Set the work directory of the image
WORKDIR /app

# Download the necessary Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the source code
# NOTE: Include files to ignore in .dockerignore
COPY . .

# Expose the port to the container
# NOTE: This is not the same port to access the application through the localhost
EXPOSE 8000

# Run the server
CMD ["go", "run", "cmd/server/main.go"]