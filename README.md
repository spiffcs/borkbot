# BorkBot

Borkbot - Goodboi for slack

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

``` bash
docker
```

### Installing

```bash
git clone git@github.com:sparklycb/borkbot.git

mkdir certs

# Generate a cert.pem and key.pem and place them in the cert folder. These are used for local ssl.
# One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
# https://golang.org/pkg/net/http/#ListenAndServeTLS
# https://gist.github.com/denji/12b3a568f092ab951456

docker build -t borkbot:dev -f ./development/Dockerfile .

docker run --rm -it -p 8080:8080 -v $(pwd):/go/src/github.com/sparklycb/borkbot borkbot:dev CompileDaemon -build="go build -o borkbotd borkbot/cmd/borkd/main.go" \
                                                                                        -command="./borkbotd --verification_token=<SLACK_VERIRICATION_TOKEN>" \
                                                                                        -exclude-dir="vendor"
```

After running the above you should see:

```bash
2018/04/03 18:09:38 Running build command!
2018/04/03 18:09:38 Build ok.
2018/04/03 18:09:38 Restarting the given command.
2018/04/03 18:09:38 stdout: ts=2018-04-03T18:09:38.8166853Z transport=https address=:9000 msg=listening
```

The health check endpoint is GET at /borkbot/v1/health which returns a 200 json response with a string.
You can test this first in postmant to make sure the bot is running.

The current application route available is a POST request at /borkbot/v1/bork which returns a random goodboi gif.

Please note: borkd requires a verificationToken to verify that requests are coming from slack.

You can test this app by setting up a slack app at: [api.slack.com](https://api.slack.com/apps/A8ZB6FMQD)

If you want to integrate with slack and see how it works when configured with your own custom app follow this guide:
https://api.slack.com/tutorials/tunneling-with-ngrok

If you don't want to take the time to setup the above slack integration you can test the server with [Postman](https://www.getpostman.com/)

Make sure your request is formatted as a POST request with the parameters listed in transport.go:74 under SlackForm
The parameters should be application/x-www-form-urlencoded where TOKEN is the SLACK_VERIFICATION_TOKEN passed in at runtime

If no verification_token is passed then a POST request to /borkbot/v1/bork will generate a valid response

## Production

To build the production image run:

```bash
docker build -t borkbot:production .
docker run --rm -it -p 443:443 -v $(pwd)/secure:/secure borkbot:production /borkbotd --listen=:443
```

The production binary can be configured with a listen flag for operator needs. Make sure your -p flag matches the listen flag specified. 8080 is the current default: cmd/borkd/main.go:17

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
