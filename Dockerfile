FROM golang:1.10-alpine3.7 as builder
# Add curl and git dependencies
RUN apk add --no-cache curl git
# Pull down godep
RUN curl -fsSL -o /usr/local/bin/dep \
        https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
        chmod +x /usr/local/bin/dep
# Set working directory
WORKDIR /go/src/github.com/sparklycb
# configure dependencies
ADD Gopkg.lock Gopkg.lock
ADD Gopkg.toml Gopkg.toml
# install dependencies
RUN dep ensure -vendor-only
# add source code
ADD . borkbot
# build the source
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o borkbotd borkbot/cmd/borkd/main.go
# use a minimal alpine image
FROM scratch
# set working directory
WORKDIR /
# copy the binary from builder
COPY --from=builder /go/src/github.com/sparklycb/borkbotd .
EXPOSE 443
# run the binary
CMD ["/borkbotd"]
