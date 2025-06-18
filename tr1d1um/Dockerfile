# Build stage
FROM golang:1.24 AS builder
 
WORKDIR /app
COPY . .
 
# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tr1d1um .
 
# Run stage
FROM alpine:latest
 
RUN mkdir -p /etc/tr1d1um
WORKDIR /etc/tr1d1um
 
COPY --from=builder /app/tr1d1um .
COPY --from=builder /app/tr1d1um.yaml .
 
RUN chmod +x tr1d1um && chmod 666 tr1d1um.yaml
 
CMD ["/etc/tr1d1um/tr1d1um"]
