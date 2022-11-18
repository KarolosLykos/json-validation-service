.PHONY: run lint test start-db remove-db


# Variables
container_name = json-validation-service
container_port = 5432

postgres_user = postgres
postgres_password = mysecretpassword
postgres_db = json-validation-service

start-db:
	docker run --name ${container_name} -p ${container_port}:5432 \
 			-e POSTGRES_USER=${postgres_user} \
 			-e POSTGRES_PASSWORD=${postgres_password} \
 			-e POSTGRES_DB=${postgres_db} \
 			-d postgres

remove-db:
	docker rm -f ${container_name}

run:
	go run cmd/main.go

lint:
	golangci-lint run -c .golangci.yml

test:
	LOGGER_LEVEL=error go test ./... -v -cover

local-up:
	docker compose up -d

local-down:
	docker compose down
