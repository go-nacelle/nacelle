FROM circleci/golang:1.15

ENV GITHUB_CLI_VERSION="1.0.0"

RUN set -ex && \
    curl -L -o gh_cli.tar.gz "https://github.com/cli/cli/releases/download/v${GITHUB_CLI_VERSION}/gh_${VERSION}_linux_amd64.tar.gz" && \
    tar -xvzf gh_cli.tar.gz "gh_${VERSION}_linux_amd64/bin/gh" && \
    mv "gh_${VERSION}_linux_amd64/bin/hub" /usr/local/bin/gh && \
    rm -rf gh_cli.tar.gz "gh_${VERSION}_linux_amd64"

# TODO - also update update.go

