# syntax=docker/dockerfile:1

################################################################################
# Create a stage for building the application.
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies as a separate step to take advantage of Docker's caching.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# This is the architecture you're building for, which is passed in by the builder.
# Placing it here allows the previous steps to be cached across architectures.
ARG TARGETARCH

# Build the application.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -ldflags="-w -s" -a -installsuffix cgo -o main .

################################################################################
# Create a new stage for running the application that contains the minimal
# runtime dependencies for the application.
FROM alpine:3.19

# Install any runtime dependencies that are needed to run your application.
RUN apk --no-cache add ca-certificates tzdata curl && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy the executable from the "build" stage.
COPY --from=builder /app/main .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app
USER appuser

# Expose the port that the application listens on.
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
    CMD curl -f http://localhost:8080/v1/health || exit 1

# What the container should run when it is started.
CMD ["./main"]
