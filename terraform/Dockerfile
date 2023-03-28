ARG BASE_IMAGE="codercom/enterprise-base:ubuntu"
FROM ${BASE_IMAGE}

USER root

# Prep
RUN apt clean \
    && apt update -y

# Install Golang
RUN curl -L "https://go.dev/dl/$(curl -Ls https://go.dev/VERSION?m=text).linux-amd64.tar.gz" \
  | tar -C /usr/local -xzv

RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile

# Install NodeJS
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash - \
    && sudo apt install -y nodejs

# Install Corepack / Yarn
RUN corepack enable \
    && npm install --global corepack

# Install ffmpeg
RUN apt install -y ffmpeg

# Install tools
RUN env GOBIN=/bin /usr/local/go/bin/go install github.com/go-task/task/v3/cmd/task@latest
RUN env GOBIN=/bin /usr/local/go/bin/go install github.com/vektra/mockery/v2@latest
RUN env GOBIN=/bin /usr/local/go/bin/go install github.com/swaggo/swag/cmd/swag@latest
RUN env GOBIN=/bin /usr/local/go/bin/go install github.com/zekrotja/schnittstelle/cmd/schnittstelle@latest
RUN env GOBIN=/bin /usr/local/go/bin/go install github.com/pressly/goose/v3/cmd/goose@latest
RUN npm install -g swagger-markdown

USER coder

RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
