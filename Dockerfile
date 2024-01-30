FROM golang:1.21.6-bookworm as golang

ARG UID=1000
ARG GID=1000

RUN groupadd -g ${GID} user && \
  useradd -ms /bin/bash -l -u ${UID} -g ${GID} user


FROM golang

WORKDIR /app
RUN chown user:user /app

COPY --chown=user:user go.mod .
COPY --chown=user:user go.sum .
USER user
RUN go mod download

USER root
COPY --chown=user:user . .
USER user
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ARG TS
RUN echo TEST ${TS} | tee -a /app/test.log && \
  go test -v ./.../ | tee -a /app/test.log
# RUN go test -v -failfast ./.../ -run SqlTemplate | tee -a /app/test.log
