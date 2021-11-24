FROM golang:1.17.3-bullseye

RUN groupadd vscode \
  && useradd -s /bin/bash --gid vscode -m vscode \
  && usermod -aG sudo vscode


RUN go get -v golang.org/x/tools/gopls