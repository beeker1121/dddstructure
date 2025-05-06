# Use the official Go image to create a binary.
FROM golang:1.23-bookworm AS builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY ../../../go.* ./
RUN go mod download

# Copy local code to the container image.
# Since we're in cmd/deployments/kubemysql, we need
# to go back three folders, ie ../../../
COPY ../../../ ./

# Build the binary.
RUN go build -v -o dddstructure cmd/api/main.go

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Create app user and group
# RUN addgroup app && adduser -S -H -D -G app app
RUN useradd -M -s /sbin/nologin app
RUN usermod -aG app app

# Copy the binary and config to the image from the builder stage.
COPY --from=builder /app/dddstructure /app/dddstructure
COPY --from=builder /app/cmd/api/config.json /app/config.json

# Give app user permissions to run dddstructure app.
RUN chown -R app /app/dddstructure

# Expose port 8080 on container (the port it listens on).
# This is more metadata, it does not automatically publish
# the port on the host.
EXPOSE 8080

# Switch to app user.
USER app

# Change to app directory
WORKDIR /app

# Run the app on container startup.
CMD ["./dddstructure"]