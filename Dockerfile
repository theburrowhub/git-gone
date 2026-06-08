# Dockerfile for git-gone
# Used by GoReleaser to build multi-arch container images

FROM alpine:3.19

# Install git (required for git-gone to work)
RUN apk add --no-cache git ca-certificates

# Create non-root user
RUN adduser -D -u 1000 gitgone

# Copy binary from GoReleaser
COPY git-gone /usr/local/bin/git-gone

# Set ownership and permissions
RUN chmod +x /usr/local/bin/git-gone

# Switch to non-root user
USER gitgone

# Set working directory
WORKDIR /repo

ENTRYPOINT ["git-gone"]

