FROM golang:1.21.6-bookworm

ARG UID=1000
ARG GID=1000

RUN groupadd -g ${GID} user && \
  useradd -ms /bin/bash -l -u ${UID} -g ${GID} user

WORKDIR /app
RUN chown user:user /app && \
  GOBIN=/usr/bin/ go install github.com/a-h/templ/cmd/templ@latest

USER user
