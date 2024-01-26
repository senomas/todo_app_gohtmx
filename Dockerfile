FROM golang:1.21.6-bookworm as golang

ARG UID=1000
ARG GID=1000

RUN groupadd -g ${GID} user && \
  useradd -ms /bin/bash -l -u ${UID} -g ${GID} user

USER user

FROM golang

WORKDIR /app

COPY --chown=user:user go.mod .
COPY --chown=user:user go.sum .
RUN go mod download

COPY --chown=user:user . .
RUN go test -v ./...
