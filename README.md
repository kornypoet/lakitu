# Lakitu

![lakitu](doc/lakitu.png)

Lakitu is a simple HTTP server written in Golang.
This repo was designed to demonstrate coding and automation fundamentals.

## Requirements

The following requirements are expected:

* golang 1.17
* make
* docker

Additionally, this document assumes the user to be running on MacOS.

## Usage

To build the binary locally for Mac, clone this repo then run `make`.

```
$ git clone https://github.com/kornypoet/lakitu.git
$ cd lakitu
$ make
```

The binary will be created in the `bin` directory. It has several different flags:

```
$ bin/darwin_amd64/lakitu -h
Simple HTTP Server

Usage:
  lakitu [flags]

Flags:
  -a, --assets string   Directory to store assets in
  -b, --bind string     Address to bind server to (default "localhost")
  -d, --debug           Enable debug mode
  -h, --help            help for lakitu
  -l, --logging         Enable request logs (default true)
  -p, --port string     Port to listen on (default "8080")
  -v, --version         version for lakitu
```

To run locally in debug mode:

```
$ bin/darwin_amd64/lakitu -d
```

This will start the server on localhost on port 8080 with additional logging and
store any downloaded assets in an `assets` directory in the repo (which has been git-ignored).
The storage location for assets can be specified with the `-a` flag, and lakitu will create the directory if it does not exist.

To check the compiled version:

```
$ bin/darwin_amd64/lakitu -v
lakitu version v0.1.0
```

This repo comes with a multi-stage `Dockerfile` for containerization. Lakitu is built with an official golang linux image
then run on an alpine linux image to keep the size small.

```
$ make docker
```

The resulting image can be run locally with docker compose. Docker compose will expose lakitu on port 9090 on localhost
and will attach the relative `assets` directory to the container for debugging and/or re-use purposes.

```
docker compose up
```

## API

Lakitu's two routes are fairly straightforward: one for retrieving the version (formerly ping) and one to manage files.

### Version

```
curl -X GET http://localhost:8080/v1/version
-> 200 vX.X.X
```

### Manage File

Manage file requires that you post a JSON payload with an action to trigger behavior in the server.
The initial implementation here made the following assumptions, and treated these calls more like RPC than REST:

* You should only download the file once; in it's current form, the API's operation is simple, but in a production setting
  this could be more complex and error-prone, so don't download the file a second time.
  Instead, return a meaningful error to the caller.

* Secondly, since we only want to download the file once, we should handle multiple concurrent calls gracefully.
  Concurrent calls are "locked" and return a meaningful error to the caller.

```
curl -X POST http://localhost:8080/v1/manage_file -d '{"action":"download"}'
# first api call
-> 200 {"action":"download","status":"success"}
# concurrent calls
-> 429 {"err":"file download in progress","status":"failure"}
# additional calls
-> 500 {"err":"file already downloaded","status":"failure"}
```

When reading the file from the server, if it hasn't been downloaded yet, we return a meaningful error to the user.
It would have been possible to implement an on-missing-download call in the API, but ultimately the decision was made to be explicit.
The file is returned as plain text, rather than as JSON

```
POST curl -X POST http://localhost:8080/v1/manage_file -d '{"action":"read"}'
# before downloading
-> 500 {"err":"file must be downloaded first","status":"failure"}
# after downloading
-> 200 Lorem ipsum ...
```

## Development

Almost all directives are defined in this repo's `Makefile`:

Build for mac (the default directive):

```
make mac
```

Build for linux (used in docker)

```
make build
```

Install tools (additional cli tools, not required for runtime)

```
make tools
```

Run tests (written using ginkgo)

```
make test
```

Format code (only outputs if changes are needed)

```
make format
```

Build docker container (also tags with the current version)

```
make docker
```

## Automation

Lakitu includes two github actions workflows. The first runs on every pull request `ci.yaml`.
It is designed to check that all tests pass and no go format changes are needed.
This check runs only on pull requests that alter go code files, the Makefile, or dependencies.

The second workflow only runs on a merge to the main branch.
It will create a new version tag when the `VERSION` file has been changed and the branch being merged has been labeled `release`.
