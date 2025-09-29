# Start from Go 1.23 base image
FROM golang:1.23 AS builder

# Install git (needed for cloning)
RUN apt-get update && apt-get install -y --no-install-recommends git \
    && rm -rf /var/lib/apt/lists/*

# Copy version.txt into the container
COPY version.txt /tmp/version.txt

# Set repo URL as build argument
ARG REPO_URL="https://github.com/rdkcentral/xconfui.git"

WORKDIR /app

# Clone based on version from version.txt
RUN VERSION=$(cat /tmp/version.txt | tr -d '[:space:]') \
    && echo "Cloning tag: $VERSION" \
    && git clone --branch "$VERSION" --depth 1 "$REPO_URL" /app

RUN git branch


# Set environment variables
ENV NVM_DIR=/root/.nvm
ENV NODE_VERSION=24.1.0
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
#WORKDIR /app

# Run npm install and grunt
RUN bash -c "source $NVM_DIR/nvm.sh && \
             nvm use $NODE_VERSION && \
             npm install && \
             npm install -g bower && \
             bower install && \
             grunt install" 

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make

# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make

FROM ubuntu:24.04

RUN mkdir -p  /app/logs/xconfadminui && mkdir -p /templates

#COPY --from=builder /app .

WORKDIR /app

COPY --from=builder /app .
 
COPY --from=builder /app/bin/xconfui-linux-amd64 .
COPY --from=builder /app/config/sample_xconfui.conf ./sample_xconfui.conf

 
RUN chmod +x xconfui-linux-amd64 && chmod 666 ./sample_xconfui.conf

RUN sed -i 's/localhost/xconfadmin/g' ./sample_xconfui.conf

EXPOSE 9090
 
CMD ["./xconfui-linux-amd64","-f","sample_xconfui.conf"]
