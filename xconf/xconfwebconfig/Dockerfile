# Start from Go 1.23 base image
FROM golang:1.23 AS builder

# Install git (needed for cloning)
RUN apt-get update && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

# Copy version.txt into the container
COPY version.txt /tmp/version.txt

# Set repo URL as build argument
ARG REPO_URL="https://github.com/rdkcentral/xconfwebconfig.git"

WORKDIR /app

# Clone based on version from version.txt
RUN VERSION=$(cat /tmp/version.txt | tr -d '[:space:]') \
    && echo "Cloning tag: $VERSION" \
    && git clone --branch "$VERSION" --depth 1 "$REPO_URL" /app

RUN git branch


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

#RUN sed -i '/^hosts[[:space:]]*=/,/]/ s/"127\.0\.0\.1"/"cassandra"/' ./xconfwebconfig.conf

RUN sed -i 's/127\.0\.0\.1/cassandra/g'  ./xconfwebconfig.conf
 
RUN chmod +x xconfwebconfig-linux-amd64 && chmod 666 ./xconfwebconfig.conf

EXPOSE 9000

CMD sh -c "/app/xconfwebconfig/xconfwebconfig-linux-amd64 & tail -f /dev/null"
