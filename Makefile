default: lintf

DB_NAME ?= runtime
DB_USER ?= `whoami`
DB_HOSTNAME ?= localhost
RD_NAME ?= runtime
RD_USER ?= `whoami`
RD_HOSTNAME ?= localhost
DB_URL="postgresql://$(DB_USER)@$(DB_HOSTNAME):5432/$(DB_NAME)"
RD_URL="redis://$(RD_USER)@$(RD_HOSTNAME):6379/$(RD_NAME)?pool=$(RD_POOL)&process=$(RD_NUM_PROCESSES)"

install:
	brew install protobuf
	brew install prototool
	brew install postgresql go
	brew install redis go
	brew services start postgresql
	brew services start redis

lintf:
	prototool lint ./protos
	prototool format -w

build:
	prototool all
	go mod tidy

test:
	dropdb --if-exists $(DB_NAME)
	createdb $(DB_NAME)
	go test -v ./... -args -db_url=$(DB_URL) -rd_url=$(RD_URL)
	dropdb $(DB_NAME)