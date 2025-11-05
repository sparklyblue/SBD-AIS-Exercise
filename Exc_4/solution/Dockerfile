# ---------- BUILD STAGE ----------
FROM golang:1.25 AS builder

WORKDIR /app

# Copy everything into /app
COPY . .

# Build the application
RUN bash scripts/build-application.sh

# ---------- RUN STAGE ----------
FROM alpine:latest

# Install needed libraries for Go binary
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /app/ordersystem .

# Expose port
EXPOSE 3000

# Run Orderservice
CMD ["./ordersystem"]
