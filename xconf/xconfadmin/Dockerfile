# Start from Go 1.23 base image
FROM golang:1.23 AS builder

# Install git (needed for cloning)
RUN apt-get update && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

# Copy version.txt into the container
COPY version.txt /tmp/version.txt

# Set repo URL as build argument
ARG REPO_URL="https://github.com/rdkcentral/xconfadmin.git"

WORKDIR /app

# Clone based on version from version.txt
RUN VERSION=$(cat /tmp/version.txt | tr -d '[:space:]') \
    && echo "Cloning tag: $VERSION" \
    && git clone --branch "$VERSION" --depth 1 "$REPO_URL" /app


# Copy credentials file into repo
COPY credentials.txt /app

RUN ADMIN_FILE=$(find /app -type f -name "auth_handler.go" | head -n1) && \
    if [ -n "$ADMIN_FILE" ]; then \
      USER=$(head -n1 /app/credentials.txt) && \
      PASS=$(tail -n1 /app/credentials.txt) && \
      sed -i "s|if authRequest\.Username.*|if authRequest.Username == \"${USER}\" \&\& authRequest.Password == \"${PASS}\" {|" "$ADMIN_FILE" ; \
      echo "Patched $ADMIN_FILE with credentials"; \
    else \
      echo 'admin.go not found' ; \
    fi

 
# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make
 
# Run stage
FROM alpine:latest
 
RUN mkdir -p /app/xconfadmin
WORKDIR /app/xconfadmin
 
COPY --from=builder /app/bin/xconfadmin-linux-amd64 .
COPY --from=builder /app/config/sample_xconfadmin.conf ./xconfadmin.conf

RUN sed -i 's/127\.0\.0\.1/cassandra/g' ./xconfadmin.conf
 
RUN chmod +x xconfadmin-linux-amd64 && chmod 666 ./xconfadmin.conf

EXPOSE 9001
 

CMD sh -c "/app/xconfadmin/xconfadmin-linux-amd64  & tail -f /dev/null"

