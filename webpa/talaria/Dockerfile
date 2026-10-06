# Build stage
ARG GO_VERSION

FROM golang:${GO_VERSION} AS builder
 
WORKDIR /app


RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*


# Define the service name (hardcoded for this Dockerfile)
ENV SERVICE=talaria
ENV REPO_URL=https://github.com/xmidt-org/talaria.git 

# Copy versions file into image
COPY versions.txt /tmp/versions.txt

# Lookup version from versions.txt and clone
RUN TAG=$(grep "^${SERVICE}-" /tmp/versions.txt | cut -d'-' -f2) && \
    if [ -z "$TAG" ]; then echo "No version found for ${SERVICE}" && exit 1; fi && \
    git clone --branch "$TAG" --depth 1 "$REPO_URL" /app


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
