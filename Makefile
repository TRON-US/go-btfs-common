default: lintf

DB_NAME ?= runtime
DB_USER ?= `whoami`
DB_HOSTNAME ?= localhost
DB_URL="postgresql://$(DB_USER)@$(DB_HOSTNAME):5432/$(DB_NAME)"

install:
	brew install protobuf
	brew install prototool
	brew install postgresql go
	brew services start postgresql

lintf:
	prototool lint ./protos
	prototool format -w

build:
	prototool all
	go mod tidy

test:
	dropdb --if-exists $(DB_NAME)
	createdb $(DB_NAME)
	go test -v ./... -args -db_url=$(DB_URL)
	dropdb $(DB_NAME)