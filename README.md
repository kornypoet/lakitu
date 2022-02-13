# Lakitu

![lakitu](doc/lakitu.png)

Lakitu is a simple HTTP server

## API

Routes

```
GET /v1/version
-> 200 vX.X.X
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

Build and run with docker compose

```
docker compose up
```

Check Version

```
bin/darwin_amd64/lakitu -v
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

Build docker container

```
make docker
```

## Requirements

* golang 1.17
* make
* docker
