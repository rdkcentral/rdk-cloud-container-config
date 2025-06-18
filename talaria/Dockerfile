# Build stage
FROM golang:1.24 AS builder
 
WORKDIR /app
COPY . .
 
# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o talaria .
 
# Run stage
FROM alpine:latest
 
RUN mkdir -p /etc/talaria
WORKDIR /etc/talaria
 
COPY --from=builder /app/talaria .
COPY --from=builder /app/talaria.yaml .
 
RUN chmod +x talaria && chmod 666 talaria.yaml
 
CMD ["/etc/talaria/talaria"]
