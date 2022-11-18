
<h1 align="center">Json Validation Service</h1>

---

<p align="center">

<a style="text-decoration: none" href="https://github.com/KarolosLykos/json-validation-service/actions?query=workflow%3ABuild+branch%3Amain">
<img src="https://img.shields.io/github/workflow/status/KarolosLykos/json-validation-service/Build?style=flat-square" alt="Build Status">
</a>

<a style="text-decoration: none" href="go.mod">
<img src="https://img.shields.io/github/go-mod/go-version/KarolosLykos/json-validation-service?style=flat-square" alt="Go version">
</a>

<a href="https://codecov.io/gh/KarolosLykos/json-validation-service" style="text-decoration: none">
<img src="https://img.shields.io/codecov/c/gh/KarolosLykos/json-validation-service?color=magenta&logo=codecov&style=flat-square" alt="Downloads">
</a>

---

## Requirements

- Go 1.17+

## Application structure

- The `cmd` folder contains the application that glues everything together.
- The `internal` folder contains the interfaces on how to communicate with the application
  - The `logruslog` folder contains a `logrus` implementation of `Logger`.
  - The `server` folder contains a `gorilla/mux` Go API server to handle HTTP transport.
  - The `store` folder contains a `postgres` implementation.
  - The `validator` folder contains the implementation of the "validating service" using `jsonschema` lib


## Dependencies

- [github.com/kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig) Manage app config from env
- [github.com/gorilla/mux](https://github.com/gorilla/mux) HTTP router (port/web)
- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus) Structured logger
- [github.com/santhosh-tekuri/jsonschema](https://github.com/santhosh-tekuri/jsonschema) Json Schema validating library
- [gorm.io/gorm](https://gorm.io/) ORM library
- [github.com/golang/mock](https://github.com/golang/mock) Mocking framework
- [github.com/stretchr/testify](https://github.com/stretchr/testify) Testing Library

## How to run this?

### Locally (requires postgres running)
```bash
DB_HOST= /
DB_PORT= /
DB_USER= /
DB_NAME= /
DB_PASSWORD= /
go run cmd/main.go
```

### Locally (with Docker)
```bash
make start-db && make run
```

The service should be accessible at <http://0.0.0.0:8082>.

### Locally (with Docker Compose)
```bash
make local-up
```

The service should be accessible at <http://0.0.0.0:8082>.


## How to run tests?

```bash
make test
```

## Endpoints

- `POST /schema/{schemaID}`

#### Example request:
```bash
curl -X POST http://localhost:8082/schema/config-schema -d @testdata/config-schema.json
```

#### Example response:
```
201 Status Created

{"action":"uploadSchema","id":"config-schema","status":"success"}
```

- `Get /schema/{schemaID}`

#### Example request:
```bash
curl -X GET http://localhost:8082/schema/config-schema
```

#### Example response:
```
200 Status OK

{"action":"downloadSchema","id":"config-schema","status":"success","payload":"{\"type\": \"object\", \"$schema\": \"http://json-schema.org/draft-04/schema#\", \"required\": [\"source\", \"destination\"], \"properties\": {\"chunks\": {\"type\": \"object\", \"required\": [\"size\"], \"properties\": {\"size\": {\"type\": \"integer\"}, \"number\": {\"type\": \"integer\"}}}, \"source\": {\"type\": \"string\"}, \"timeout\": {\"type\": \"integer\", \"maximum\": 32767, \"minimum\": 0}, \"destination\": {\"type\": \"string\"}}}"}
```

#### Example request:
```bash
curl -X POST http://localhost:8082/validate/config-schema -d @testdata/config.json
```

#### Example response:
```
200 Status OK

{"action":"validateSchema","id":"config-schema","status":"success"}
```

## Extras
- Basic unit test on `Upload, Download and Validate handlers` and `validator service`
- Added `Github actions` for linting, testing and building the service.