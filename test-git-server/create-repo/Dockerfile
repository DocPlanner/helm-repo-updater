FROM alpine:3.15.0

LABEL maintainer="Docplanner <devops@docplanner.com>"

RUN apk add --no-cache bash git

COPY create-repo/start.sh .

CMD ["bash", "start.sh"]