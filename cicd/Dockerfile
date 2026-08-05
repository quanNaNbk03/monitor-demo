FROM golang:1.25.0-bookworm AS go-builder

RUN apt-get update && apt-get install -y \
    git \
    make \
    nodejs \
    npm \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
ENV TZ=Asia/Ho_Chi_Minh
ENV GOPRIVATE=git.ocn.com.vn
ENV GOINSECURE=$GOPRIVATE

ARG GIT_TOKEN
ARG GIT_USER

RUN test -n "${GIT_TOKEN}" || (echo "ERROR: GIT_TOKEN is not set" && exit 1)
RUN test -n "${GIT_USER}" || (echo "ERROR: GIT_USER is not set" && exit 1)

RUN git config --global http.sslVerify false && \
    git config --global url."https://git.ocn.com.vn/".insteadOf "git@git.ocn.com.vn:"

RUN echo "machine git.ocn.com.vn" > /root/.netrc && \
    echo "login ${GIT_USER}" >> /root/.netrc && \
    echo "password ${GIT_TOKEN}" >> /root/.netrc && \
    chmod 600 /root/.netrc

RUN go env -w GOPRIVATE=$GOPRIVATE GOINSECURE=$GOINSECURE

COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

# Copy source code
COPY . .

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Run with verbose output
RUN set -x && make generate && make build

# ################################################
# STAGE 2.1: Build deb file using binary for Debian 12 
# ################################################
# Use the official Debian 12 (bookworm) base image
FROM debian:12-slim AS deb-builder-debian12

WORKDIR /app

# Args
ARG NEXUS_PASSWORD

# Install dependencies for Debian 12
RUN ln -sf /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime

# Configure Debian repositories for better package availability
RUN echo "deb http://deb.debian.org/debian bookworm main" > /etc/apt/sources.list && \
    echo "deb http://deb.debian.org/debian bookworm-updates main" >> /etc/apt/sources.list && \
    echo "deb http://security.debian.org/debian-security bookworm-security main" >> /etc/apt/sources.list

# Update and install build tools
RUN apt update && \
    apt install -y --no-install-recommends \
    build-essential \
    dh-make \
    debhelper \
    dpkg-dev \
    devscripts \
    tree \
    openssh-client \
    libterm-readkey-perl \
    curl \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy odns-api binary application
COPY --from=go-builder /app/ /app/
RUN rm -f ./odns-api-admin/opt/ocn/odns-api-admin/bin/.gitkeep
RUN rm -f ./odns-api-client/opt/ocn/odns-api-client/bin/.gitkeep
RUN cp ./odns-api-admin/opt/ocn/odns-api-admin/bin/odns-api-admin ./odns-api-client/opt/ocn/odns-api-client/bin/odns-api-client

# Verify all data folder and binary
RUN echo "=== VERIFYING PACKAGE STRUCTURE ===" && \
    cp ./api/openapi.bundled.json ./odns-api-admin/opt/ocn/odns-api-admin/config/ && \
    cp ./api/openapi.bundled.json ./odns-api-client/opt/ocn/odns-api-client/config/ && \
    tree ./odns-api-admin && \
    tree ./odns-api-client && \
    echo "=== CHECKING BINARY ===" && \
    ls -lah ./odns-api-admin/opt/ocn/odns-api-admin/bin/ && \
    ls -lah ./odns-api-client/opt/ocn/odns-api-client/bin/ && \
    if [ -f "./odns-api-admin/opt/ocn/odns-api-admin/bin/odns-api-admin" ]; then \
    echo "✅ Binary found"; \
    ls -lah ./odns-api-admin/opt/ocn/odns-api-admin/bin/odns-api-admin; \
    echo "Binary size: $(stat -c%s ./odns-api-admin/opt/ocn/odns-api-admin/bin/odns-api-admin) bytes"; \
    else \
    echo "❌ Binary NOT found - checking what was copied:"; \
    find ./odns-api-admin -name "odns-api-admin" -type f; \
    exit 1; \
    fi

# Build deb file with cache invalidation
RUN echo "Build timestamp: $(date)" && \
    find ./odns-api-admin/DEBIAN/ ! -name control -exec chmod 0555 {} \; && chmod 0755 ./odns-api-admin/DEBIAN
RUN dpkg-deb --build odns-api-admin && ls -lah ./ && \
    echo "=== DEB FILE VERIFICATION ===" && \
    ls -lah ./odns-api-admin.deb && \
    file ./odns-api-admin.deb && \
    dpkg -I ./odns-api-admin.deb

RUN echo "Build timestamp: $(date)" && \
    find ./odns-api-client/DEBIAN/ ! -name control -exec chmod 0555 {} \; && chmod 0755 ./odns-api-client/DEBIAN
RUN dpkg-deb --build odns-api-client && ls -lah ./ && \
    echo "=== DEB FILE VERIFICATION ===" && \
    ls -lah ./odns-api-client.deb && \
    file ./odns-api-client.deb && \
    dpkg -I ./odns-api-client.deb

# Upload repo Debian 12 (bookworm)
RUN curl -u "cicd:${NEXUS_PASSWORD}" -fail -k -X POST \
    "https://nexus.ocn.com.vn/service/rest/v1/components?repository=ubuntu-jammy" \
    -H "accept: application/json" \
    -H "Content-Type: multipart/form-data" \
    -F "apt.asset=@./odns-api-admin.deb"

RUN curl -u "cicd:${NEXUS_PASSWORD}" -fail -k -X POST \
    "https://nexus.ocn.com.vn/service/rest/v1/components?repository=ubuntu-jammy" \
    -H "accept: application/json" \
    -H "Content-Type: multipart/form-data" \
    -F "apt.asset=@./odns-api-client.deb"
