# Build stage
FROM golang:1.24 AS builder
 
WORKDIR /app
COPY . .
 
# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scytale .
 
# Run stage
FROM alpine:latest
 
RUN mkdir -p /etc/scytale
WORKDIR /etc/scytale
 
COPY --from=builder /app/scytale .
COPY --from=builder /app/scytale.yaml .
 
RUN chmod +x scytale && chmod 666 scytale.yaml
 
CMD ["/etc/scytale/scytale"]
