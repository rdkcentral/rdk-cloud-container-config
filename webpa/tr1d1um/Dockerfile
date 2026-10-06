ARG GO_VERSION

FROM golang:${GO_VERSION} AS builder

 
WORKDIR /app

RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*


# Define the service name (hardcoded for this Dockerfile)
ENV SERVICE=tr1d1um
ENV REPO_URL=https://github.com/xmidt-org/tr1d1um.git 

# Copy versions file into image
COPY versions.txt /tmp/versions.txt

# Lookup version from versions.txt and clone
RUN TAG=$(grep "^${SERVICE}-" /tmp/versions.txt | cut -d'-' -f2) && \
    if [ -z "$TAG" ]; then echo "No version found for ${SERVICE}" && exit 1; fi && \
    git clone --branch "$TAG" --depth 1 "$REPO_URL" /app


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
