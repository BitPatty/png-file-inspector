FROM golang:1.17.3-bullseye

RUN groupadd vscode \
  && useradd -s /bin/bash --gid vscode -m vscode \
  && usermod -aG sudo vscode

USER vscode


RUN GOPATH=/go cd /tmp && go get -v golang.org/x/tools/gopls \
  golang.org/x/tools/gopls \
  github.com/uudashr/gopkgs/v2/cmd/gopkgs \
  github.com/ramya-rao-a/go-outline \
  github.com/go-delve/delve/cmd/dlv \
  honnef.co/go/tools/cmd/staticcheck