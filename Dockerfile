FROM golang:1.21.6-bookworm as golang

ARG UID=1000
ARG GID=1000

RUN groupadd -g ${GID} user && \
  useradd -ms /bin/bash -l -u ${UID} -g ${GID} user


FROM golang

WORKDIR /app
RUN chown user:user /app && \
  GOBIN=/usr/bin/ go install github.com/a-h/templ/cmd/templ@latest && \
  chown -R user:user /go/pkg

USER user

COPY --chown=user:user go.mod .
COPY --chown=user:user go.sum .

RUN go mod download

COPY --chown=user:user . .

RUN templ generate

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
ARG TS
RUN echo TEST ${TS} | tee -a /app/test.log && \
  MIGRATIONS_PATH=/app/migrations \
  go test -v ./store/.../ | tee -a /app/test.log

RUN echo TEST ${TS} | tee -a /app/test.log && \
  MIGRATIONS_PATH=/app/migrations \
  go test -v ./handler/.../ | tee -a /app/test.log
