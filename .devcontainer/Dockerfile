FROM golang:1.17.5

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# This Dockerfile adds a non-root 'vscode' user with sudo access. However, for Linux,
# this user's GID/UID must match your local user UID/GID to avoid permission issues
# with bind mounts. Update USER_UID / USER_GID if yours is not 1000. See
# https://aka.ms/vscode-remote/containers/non-root-user for details.
ARG USERNAME=vscode
ARG USER_UID=1000
ARG USER_GID=$USER_UID

ENV GO111MODULE=auto

ENV GOLANGCI_LINT_V="v1.43.0"
ENV DOCKER_COMPOSE_V="1.29.2"

# Configure apt, install packages and tools
RUN apt-get update \
    && apt-get -y install --no-install-recommends apt-utils dialog 2>&1 \
    #
    # Verify git, process tools, lsb-release (common in install instructions for CLIs) installed
    && apt-get -y install git iproute2 procps lsb-release \
    #
    # Install gocode-gomod
    && go get -x -d github.com/stamblerre/gocode 2>&1 \
    && go build -o gocode-gomod github.com/stamblerre/gocode \
    && mv gocode-gomod $GOPATH/bin/ \
    #
    # Install Go tools
    && go install golang.org/x/tools/...@v0.1.8 \
    && go install github.com/uudashr/gopkgs/cmd/gopkgs@v1.3.2 \
    && go install github.com/ramya-rao-a/go-outline@v0.0.0-20210608161538-9736a4bde949 \
    && go install github.com/godoctor/godoctor@v0.0.0-20211103152449-96567efaee44 \
    && go install github.com/rogpeppe/godef@v1.1.2 \
    && go install github.com/zmb3/gogetdoc@v0.0.0-20190228002656-b37376c5da6a \
    && go install github.com/sqs/goreturns@v0.0.0-20181028201513-538ac6014518 \
    && go install github.com/josharian/impl@v1.1.0 \
    && go install github.com/davidrjenni/reftools/cmd/fillstruct@v0.0.0-20210213085015-40322ffdc2e4 \
    && go install github.com/fatih/gomodifytags@v1.16.0 \
    && go install github.com/cweill/gotests/...@v1.6.0 \
    && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0 \
    && go install honnef.co/go/tools/...@v0.2.2 \
    && go install github.com/mgechev/revive@v1.1.2 \
    # Install golangci-lint
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_LINT_V} \
    #
    # Create a non-root user to use if preferred - see https://aka.ms/vscode-remote/containers/non-root-user.
    && groupadd --gid $USER_GID $USERNAME \
    && useradd -s /bin/bash --uid $USER_UID --gid $USER_GID -m $USERNAME \
    # [Optional] Add sudo support
    && apt-get install -y sudo \
    && echo $USERNAME ALL=\(root\) NOPASSWD:ALL > /etc/sudoers.d/$USERNAME \
    && chmod 0440 /etc/sudoers.d/$USERNAME \
    # Docker install
    && apt-get install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common lsb-release \
    && curl -fsSL https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | apt-key add - 2>/dev/null \
    && add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') $(lsb_release -cs) stable" \
    && apt-get update \
    && apt-get install -y docker-ce-cli \
    # Docker-compose install
    && curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_V}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
    && chmod +x /usr/local/bin/docker-compose && ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose \
    #
    # Install pip & pre-commit
    && apt-get -y install python3-pip \
    && python3 -m pip install --no-cache-dir pre-commit \
    #
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*