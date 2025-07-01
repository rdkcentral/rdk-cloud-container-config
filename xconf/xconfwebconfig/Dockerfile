# Build stage
FROM golang:1.24 AS builder
 
WORKDIR /app

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

RUN git clone https://github.com/rdkcentral/xconfwebconfig.git .

COPY . .

# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make
 
# Run stage
FROM alpine:latest
 
RUN mkdir -p /app/xconfwebconfig
RUN mkdir -p /app/logs/xconfwebconfig
RUN touch /app/logs/xconfwebconfig/xconfwebconfig.log

WORKDIR /app/xconfwebconfig
 
COPY --from=builder /app/bin/xconfwebconfig-linux-amd64  .
COPY --from=builder /app/config/sample_xconfwebconfig.conf ./xconfwebconfig.conf
 
RUN chmod +x xconfwebconfig-linux-amd64 && chmod 666 ./xconfwebconfig.conf

EXPOSE 9000

CMD sh -c "/app/xconfwebconfig/xconfwebconfig-linux-amd64 & tail -f /dev/null"
