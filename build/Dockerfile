FROM golang:1.23.2 AS build

ENV BIN_FILE=/opt/antibruteforce/ratelimiter
ENV CODE_DIR=/go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/ratelimiter/*

FROM alpine:3.9

ENV BIN_FILE="/opt/antibruteforce/ratelimiter"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/antibruteforce/config.toml
COPY ./configs/config.toml ${CONFIG_FILE}

CMD ["sh", "-c", "${BIN_FILE}"]
