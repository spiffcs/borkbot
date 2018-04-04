FROM golang:1.10-alpine3.7

RUN apk add --no-cache curl git

RUN curl -fsSL -o /usr/local/bin/dep \
        https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
        chmod +x /usr/local/bin/dep

RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /go/src/github.com/sparklycb

ADD Gopkg.lock Gopkg.lock
ADD Gopkg.toml Gopkg.toml

RUN dep ensure -vendor-only

ADD . borkbot

EXPOSE 9000

CMD ["go", "run", "src/cmd/borkd/main.go"]