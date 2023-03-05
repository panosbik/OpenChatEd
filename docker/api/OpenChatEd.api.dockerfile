ARG GO_VERSION=1.19

# Use the official Golang image as the base image
FROM golang:${GO_VERSION} AS builder

# Install required system packages
RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg openssl gcc libc6-dev libgeos-dev dumb-init

# Set the working directory to /www/src/OpenChatEd
WORKDIR /www/src/OpenChatEd/api

# Copy the Go module files and download the dependencies
COPY . .

# Download the dependencies
RUN go mod download

# Copy the private and bublic key for token authorization 
COPY ./auth.cert ./auth.cert

# Build the application
RUN go build -o .

# Use dumb-init as the entrypoint
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# Start the application
CMD ["./OpenChatEd"]
