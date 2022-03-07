FROM --platform=$BUILDPLATFORM golang:1.17-alpine3.14 as builder
RUN apk --no-cache add git
WORKDIR /go/src/build
COPY . .

# https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide/
ARG TARGETOS
ARG TARGETARCH
RUN mkdir -p dist \
    && go mod vendor \
    && CGO_ENABLED=0 GOOS=$(TARGETOS) GOARCH=$(TARGETARCH) go build -o dist/helm-repo-updater .

FROM alpine:3.14
ENV SSH_KNOWN_HOSTS="~/.ssh/known_hosts"
RUN apk update && apk add openssh
RUN mkdir -p ~/.ssh/
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
COPY hack/ /usr/local/bin/
COPY --from=builder /go/src/build/dist/ .

ENTRYPOINT ["./helm-repo-updater"]
