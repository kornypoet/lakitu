# Lakitu

![lakitu](doc/lakitu.png)

Lakitu is a simple HTTP server

## API

Routes

```
GET /v1/ping
-> 200 pong
POST /v1/manage_file {"action":"download"}
# first api call
-> 200 {"action":"download","status":"success"}
# concurrent calls
-> 429 {"err":"file download in progress","status":"failure"}
# additional calls
-> 500 {"err":"file already downloaded","status":"failure"}
POST /v1/manage_file {"action":"read"}
# before downloading
-> 500 {"err":"file must be downloaded first","status":"failure"}
# after downloading
-> 200 Lorem ipsum ...
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

Install tools

```
make tools
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
