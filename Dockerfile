FROM golang:1.17-alpine3.14 as builder
RUN apk --no-cache add git
WORKDIR /go/src/build
COPY . .
RUN export CGO_ENABLED=0 \
    && mkdir -p dist \
    && go mod vendor \
    && go build -o dist/helm-repo-updater .

FROM alpine:3.14
RUN wget -O /usr/local/bin/yq https://github.com/mikefarah/yq/releases/download/v4.14.1/yq_linux_amd64 \
    && chmod +x /usr/local/bin/yq
COPY hack/ /usr/local/bin/
COPY --from=builder /go/src/build/dist/ .

ENTRYPOINT ["./helm-repo-updater"]