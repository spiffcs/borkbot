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

docker build -t borkbot .

docker run --rm -it -p 9000:9000 -v $(pwd):/go/src/github.com/sparklycb/borkbot borkbot sh

/go/src/github.com/sparklycb: CompileDaemon -build="go build -o borkbotd borkbot/cmd/borkd/main.go" \
                                            -command="./borkbotd --verification_token=<SLACK_VERIRICATION_TOKEN>" \
                                            -exclude-dir="vendor"
```

After running the above you should see:

```bash
2018/04/03 18:09:38 Running build command!
2018/04/03 18:09:38 Build ok.
2018/04/03 18:09:38 Restarting the given command.
2018/04/03 18:09:38 stdout: ts=2018-04-03T18:09:38.8166853Z transport=http address=:9000 msg=listening
```

The current route available is a POST request at /borkbot/v1/bork which returns a random goodboi gif.

Please note: borkd requires a verificationToken to verify that requests are coming from slack.

You can test this app by setting up a slack app at: [api.slack.com](https://api.slack.com/apps/A8ZB6FMQD)

If you want to integrate with slack and see how it works when configured with your own custom app follow this guide:
https://api.slack.com/tutorials/tunneling-with-ngrok

If you don't want to take the time to setup the above slack integration you can test the server with [Postman](https://www.getpostman.com/)
Make sure your request is formatted as a POST request with the parameters listed in transport.go:74 under SlackForm
The parameters should be application/x-www-form-urlencoded where TOKEN is the SLACK_VERIFICATION_TOKEN passed in at runtime

## Production
To build the production image run:

```bash
docker build -t borkbot:production -f ./production/Dockerfile .
docker run -d -p 9000:9000 borkbot:production borkbotd --verification_token=<SLACK_VERIRICATION_TOKEN>
```

The production binary can be configured with a listen flag for operator needs. Make sure your -p flag matches the listen flag specified. 9000 is the current default: cmd/borkd/main.go:17

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
