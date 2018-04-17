FROM golang:1.10-alpine3.7

RUN apk add --no-cache curl git

RUN curl -fsSL -o /usr/local/bin/dep \
        https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
        chmod +x /usr/local/bin/dep

RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /go/src/github.com/sparklycb

ADD Gopkg.lock Gopkg.lock
ADD Gopkg.toml Gopkg.toml
ADD secure/cert.pem /secure/cert.pem
ADD secure/key.pem /secure/key.pem

RUN dep ensure -vendor-only

ADD . borkbot

EXPOSE 8080

CMD ["go", "run", "borkbot/cmd/borkd/main.go"]
