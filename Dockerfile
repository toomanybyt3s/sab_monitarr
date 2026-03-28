FROM golang:1.26-alpine AS builder

WORKDIR /app

# Set build arguments for optimization
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Download HTMX at build time so it is never committed to the repository
RUN wget -q -O /tmp/htmx.min.js \
    https://unpkg.com/htmx.org@2.0.8/dist/htmx.min.js

# Copy only necessary files for building
COPY go.* ./
RUN go mod download

# Now copy source code and build with optimization flags
COPY main.go main_test.go ./
COPY static/ static/
COPY templates/ templates/

# Place the downloaded HTMX into the static directory
RUN cp /tmp/htmx.min.js static/htmx.min.js

RUN go build -ldflags="-s -w" -o sab_monitarr .

# Use distroless as minimal base image
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copy only the necessary files from builder
COPY --from=builder /app/sab_monitarr .
COPY --from=builder /app/templates/ templates/
COPY --from=builder /app/static/ static/

# Default environment variables
ENV SABMON_REFRESH_INTERVAL=5 \
    SABMON_DEBUG=false \
    SABMON_LOG_CLIENT_INFO=false

EXPOSE 5959

# Run as non-root
USER nonroot:nonroot
CMD ["/app/sab_monitarr"]
