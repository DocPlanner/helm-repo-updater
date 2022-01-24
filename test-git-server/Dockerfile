FROM alpine:3.15.0

LABEL maintainer="Docplanner <devops@docplanner.com>"

RUN apk add --no-cache openssh git

WORKDIR /git-server/

# -D flag avoids password generation
# -s flag changes user's shell
RUN mkdir /git-server/keys \
  && adduser -D -s /usr/bin/git-shell git \
  && echo git:12345 | chpasswd \
  && mkdir /home/git/.ssh


# Create repos folder
RUN mkdir /git-server/repos/

# This is a login shell for SSH accounts to provide restricted Git access.
# It permits execution only of server-side Git commands implementing the
# pull/push functionality, plus custom commands present in a subdirectory
# named git-shell-commands in the user’s home directory.
# More info: https://git-scm.com/docs/git-shell
COPY git-server/git-shell-commands /home/git/git-shell-commands
RUN chmod -R +x /home/git/git-shell-commands

# sshd_config file is edited for enable access key and disable access password
COPY git-server/sshd_config /etc/ssh/sshd_config
COPY git-server/create_repo.sh create_repo.sh
COPY git-server/start.sh start.sh

# copy predefined ssh public key
COPY repo_keys/id_rsa.pub /git-server/keys/

# copy predefined ssh private key
COPY private_keys/helm-repo-updater-test /etc/ssh
RUN chmod 600 /etc/ssh/helm-repo-updater-test

EXPOSE 22

CMD ["sh", "start.sh"]
