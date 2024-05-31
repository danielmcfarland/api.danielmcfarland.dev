# Micropub API

## Introduction

This application can be used to gain an `access_token` in order to create a post on [danielmcfarland.dev](https://danielmcfarland.dev).

To log in visit [api.danielmcfarland.dev/auth](https://api.danielmcfarland.dev/auth).

You will be prompted to login with IndieAuth and then be redirected back with an access token.

This token can then be used as an parameter when making a request to [api.danielmcfarland.dev/micropub](https://api.danielmcfarland.dev/micropub)

You will need a `GITHUB_API_KEY` in order to write to the repo and add the post contents.

The GitHub repo this writes to is hardcoded. It is my intention to make this more configurable in the future.

```
curl --location 'https://api.danielmcfarland.dev/micropub' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'h=entry' \
--data-urlencode 'access_token=####' \
--data-urlencode 'bookmark-of=https://danielmcfarland.dev'
```

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

An `.env` file is required with the following parameters locally

```
PORT=8080
APP_URL=https://local_url_for_callback
GITHUB_API_KEY=
```

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

## Deployment

I deploy this to fly.io, you will need to add a secret token for GI# Micropub API

## Introduction

This application can be used to gain an `access_token` in order to create a post on [danielmcfarland.dev](https://danielmcfarland.dev).

To log in visit [api.danielmcfarland.dev/auth](https://api.danielmcfarland.dev/auth).

You will be prompted to login with IndieAuth and then be redirected back with an access token.

This token can then be used as an parameter when making a request to [api.danielmcfarland.dev/micropub](https://api.danielmcfarland.dev/micropub)

You will need a `GITHUB_API_KEY` in order to write to the repo and add the post contents.

The GitHub repo this writes to is hardcoded. It is my intention to make this more configurable in the future.

```
curl --location 'https://api.danielmcfarland.dev/micropub' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'h=entry' \
--data-urlencode 'access_token=####' \
--data-urlencode 'bookmark-of=https://danielmcfarland.dev'
```

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

An `.env` file is required with the following parameters locally

```
PORT=8080
APP_URL=https://local_url_for_callback
GITHUB_API_KEY=
```

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

## Deployment

I deploy this to [fly.io](https://fly.io), you will need to add a secret token for `GITHUB_API_KEY` to the application on [fly.io](https://fly.io)