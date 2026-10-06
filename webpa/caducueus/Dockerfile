ARG GO_VERSION

FROM golang:${GO_VERSION} AS builder

 
WORKDIR /app


RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*


# Define the service name (hardcoded for this Dockerfile)
ENV SERVICE=caduceus
ENV REPO_URL=https://github.com/xmidt-org/caduceus.git 

# Copy versions file into image
COPY versions.txt /tmp/versions.txt

# Lookup version from versions.txt and clone
RUN TAG=$(grep "^${SERVICE}-" /tmp/versions.txt | cut -d'-' -f2) && \
    if [ -z "$TAG" ]; then echo "No version found for ${SERVICE}" && exit 1; fi && \
    git clone --branch "$TAG" --depth 1 "$REPO_URL" /app


# Build a static binary compatible with Alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o caduceus .
 
# Run stage
FROM alpine:latest
 
RUN mkdir -p /etc/caduceus
WORKDIR /etc/caduceus
 
COPY --from=builder /app/caduceus .
COPY --from=builder /app/caduceus.yaml .
 
RUN chmod +x caduceus && chmod 666 caduceus.yaml
 
CMD ["/etc/caduceus/caduceus"]
