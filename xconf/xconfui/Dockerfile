# Build stage
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make


FROM ubuntu:24.04

# Set environment variables
ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION=20.17.0
ENV PATH=$NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

# Install dependencies
RUN apt-get update && \
    apt-get install -y curl ca-certificates build-essential git && \
    rm -rf /var/lib/apt/lists/*

# Install NVM
RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# Install Node.js and Grunt CLI
RUN bash -c "source $NVM_DIR/nvm.sh && \
             nvm install $NODE_VERSION && \
             nvm use $NODE_VERSION && \
             nvm alias default $NODE_VERSION && \
             npm install -g npm grunt-cli"

# Copy your project files (adjust as needed)
WORKDIR /app
COPY . .

# Run npm install and grunt
RUN bash -c "source $NVM_DIR/nvm.sh && \
             nvm use $NODE_VERSION && \
             npm install && \
             grunt install"

 
COPY --from=builder /app/bin/xconfui-linux-amd64 .
COPY --from=builder /app/config/sample_xconfui.conf ./sample_xconfui.conf
 
RUN chmod +x xconfui-linux-amd64 && chmod 666 sample_xconfui.conf

EXPOSE 9090
 
CMD ["./xconfui-linux-amd64","-f","sample_xconfui.conf"]
