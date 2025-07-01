# Build stage
FROM golang:1.24 AS builder
 
WORKDIR /app
COPY . .
 
# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make
 
# Run stage
FROM alpine:latest
 
RUN mkdir -p app/xconfadmin
WORKDIR app/xconfadmin
 
COPY --from=builder /app/bin/xconfadmin-linux-amd64 .
COPY --from=builder /app/config/sample_xconfadmin.conf ./xconfwebconfig.conf
 
RUN chmod +x xconfadmin-linux-amd64 && chmod 666 ./xconfwebconfig.conf

EXPOSE 9001
 

CMD sh -c "/app/xconfadmin/xconfadmin-linux-amd64  & tail -f /dev/null"

