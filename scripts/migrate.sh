#!/usr/bin/env bash

up() {
 migrate -database "${POSTGRESQL_URL}" -path database/migrations up

}

down() {
 migrate -database "${POSTGRESQL_URL}" -path database/migrations down

}

create() {
  migrate create -ext sql -dir database/migrations -seq "$1"
}

"$@"