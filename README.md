# Bullwark Microservice Server

This project is a service which aims to deploy microservices into the Bullwark
pen. Currently, only standalone JAR files are supported, which are injected into
a Docker container running OpenJDK.

## Prerequisites

- Go 1.10+ OR Docker


## Build

With native Go:

```bash
make build
```

With Docker:

```bash
make docker-build
```

## Run (Server Mode)

```bash
./bin/bullwark-microservice-server
```

You can then access the server at http://localhost:8000.
