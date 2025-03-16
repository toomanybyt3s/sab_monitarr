FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o sab_monitarr .

FROM alpine:3.21.3

WORKDIR /app
COPY --from=builder /app/sab_monitarr .
COPY templates/ /app/templates/
COPY static/ /app/static/
COPY favicon/ /app/favicon/

# Default environment variables
ENV SABMON_REFRESH_INTERVAL=5
ENV SABMON_DEBUG=false
ENV SABMON_LOG_CLIENT_INFO=false

EXPOSE 5959
VOLUME ["/app/config"]

CMD ["./sab_monitarr"]
