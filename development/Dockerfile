FROM golang:1.12-alpine3.9

RUN apk add --no-cache curl git

RUN curl -fsSL -o /usr/local/bin/dep \
        https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
        chmod +x /usr/local/bin/dep

RUN go get github.com/githubnemo/CompileDaemon

WORKDIR /go/src/github.com/sparklycb

ADD secure/cert.pem /secure/cert.pem
ADD secure/key.pem /secure/key.pem

ADD . borkbot

EXPOSE 8080

CMD ["go", "run", "borkbot/cmd/borkd/main.go"]
