# Use the official Golang image to create a build artifact.
FROM golang:1.21 AS builder

# Set the working directory inside the container.
WORKDIR /app

# Copy go mod and sum files and download dependencies
COPY src/go.mod src/go.sum ./
RUN go mod download

# Copy the entire source code from the src directory to the working directory in the container
COPY src/ .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use the scratch image for a smaller footprint.
FROM scratch
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]