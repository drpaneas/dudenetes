FROM alpine:3.7
MAINTAINER Panos Georgiadis <drpaneas@gmail.com>

COPY dudenetes /usr/local/bin/dudenetes
COPY godog /usr/local/bin/godog

ENV SHELL="/bin/sh"

# openssh-client: required for ssh-agent (for SUSE CaaSP)
# curl: required to fetch kubectl binary and the ssh key (for SUSE CaaSP)
RUN apk add --no-cache openssh-client curl git \
    && curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.15.1/bin/linux/amd64/kubectl \
    && chmod u+x kubectl && mv kubectl /bin/kubectl \
    && eval "$(ssh-agent -s)" \
    && curl -LO https://raw.githubusercontent.com/SUSE/skuba/master/ci/infra/id_shared \
    && chmod 0400 id_shared \
    && ssh-add id_shared \
    && mkdir /vmware

WORKDIR /rug

ENTRYPOINT ["/usr/local/bin/dudenetes"]
CMD ["test", "-f", "/rug", "-s", "/vmware"]
