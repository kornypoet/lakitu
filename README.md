# Lakitu

![lakitu](doc/lakitu.png)

Lakitu is a simple HTTP server

## API

Routes

```
GET /v1/ping
-> pong
POST /v1/manage_file {"action":"jackson"}
-> {"status":"success"}
```

## Usage

Build and run locally in debug mode

```
make
bin/darwin_amd64/lakitu -d
```

## Development

Build for mac

```
make mac
```

Build for linux

```
make build
```

Run tests

```
make test
```

Format code

```
make format
```

## Requirements

* golang 1.17
* make
