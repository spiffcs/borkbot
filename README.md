# Project Title

Borkbot - Goodboi for slack

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

What things you need to install the software and how to install them

``` bash
golang
dep
```

### Installing

A step by step series of examples that tell you have to get a development env running

* clone down borkbot

``` bash
git clone git@github.com:spiffcs/borkbot.git
```

* run dep to get all the dependencies from the root borkbot folder

``` bash
dep ensure
```

* make sure you can build and execute the server

``` bash
go build ./cmd/borkd/
./borkd
```

It should be listening on :9000
The route available is a get request at /bork which returns a static goodboi gif

Please note: borkd requires a verificationToken to verify that requests are coming from slack

If you want to integrate and see how it works when configured with slack follow this guide:
https://api.slack.com/tutorials/tunneling-with-ngrok

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
