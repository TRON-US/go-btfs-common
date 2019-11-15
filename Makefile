default: lintf

TEST_DB_NAME ?= runtime
TEST_DB_USER ?= `whoami`
TEST_DB_HOSTNAME ?= localhost
TEST_DB_URL="postgresql://$(TEST_DB_USER)@$(TEST_DB_HOSTNAME):5432/$(TEST_DB_NAME)"
TEST_RD_NAME ?= runtime
TEST_RD_USER ?= `whoami`
TEST_RD_HOSTNAME ?= localhost
TEST_RD_URL="redis://$(TEST_RD_USER)@$(TEST_RD_HOSTNAME):6379/$(TEST_RD_NAME)?pool=$(TEST_RD_POOL)&process=$(TEST_RD_NUM_PROCESSES)"

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
	dropdb --if-exists $(TEST_DB_NAME)
	createdb $(TEST_DB_NAME)
	go test -v ./... -args -db_url=$(TEST_DB_URL) -rd_url=$(TEST_RD_URL)
	dropdb $(TEST_DB_NAME)