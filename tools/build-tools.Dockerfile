FROM ubuntu:18.04

# Install prerequisite
RUN apt update && \
    apt-get install software-properties-common -y
RUN apt-add-repository ppa:git-core/ppa && \
    apt update && \
    apt install -y wget curl build-essential git

# Use Bash instead of Dash
RUN ln -sf bash /bin/sh

# Install docker client
RUN curl -LO https://download.docker.com/linux/static/stable/x86_64/docker-19.03.2.tgz && \
    docker_sha256=865038730c79ab48dfed1365ee7627606405c037f46c9ae17c5ec1f487da1375 && \
    echo "$docker_sha256 docker-19.03.2.tgz" | sha256sum -c - && \
    tar xvzf docker-19.03.2.tgz && \
    mv docker/* /usr/local/bin && \
    rm -rf docker docker-19.03.2.tgz

# Install golang
RUN GO_VERSION=1.17.5 && \
    curl -LO https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    go_sha256=bd78114b0d441b029c8fe0341f4910370925a4d270a6a590668840675b0c653e && \
    echo "$go_sha256 go${GO_VERSION}.linux-amd64.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xvzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm -rf go${GO_VERSION}.linux-amd64.tar.gz

ENV PATH=${PATH}:/usr/local/go/bin \
    GOROOT=/usr/local/go \
    GOPATH=/go

# Copy private key
COPY test-git-server/private_keys/helm-repo-updater-test .

RUN mkdir -p ~/.ssh
RUN chmod 600 ./helm-repo-updater-test

RUN echo "[localhost]:2222 $(ssh-keygen -f ./helm-repo-updater-test -y | cut -d' ' -f-2)" >> ~/.ssh/known_hosts