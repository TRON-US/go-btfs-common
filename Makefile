UNAME := $(shell uname)
default: lintf

TEST_DB_NAME ?= runtime
TEST_DB_USER ?= `whoami`
TEST_DB_HOSTNAME ?= localhost
TEST_DB_URL="postgresql://$(TEST_DB_USER)@$(TEST_DB_HOSTNAME):5432/$(TEST_DB_NAME)"
TEST_RD_NAME ?= runtime
TEST_RD_USER ?= `whoami`
TEST_RD_HOSTNAME ?= localhost
TEST_RD_URL="redis://$(TEST_RD_USER)@$(TEST_RD_HOSTNAME):6379/$(TEST_RD_NAME)"

PG_FIX_CANDIDATES=./protos/node/node.pb.go \
			./protos/status/status.pb.go \
			./protos/escrow/escrow.pb.go \
			./protos/guard/guard.pb.go \

install: brew trongogo

brew:
	brew install protobuf
	brew install prototool
	brew install postgresql go
	brew install redis go
	brew services start postgresql
	brew services start redis

trongogo:
	cd ../ && git clone https://github.com/TRON-US/protobuf || true
	cd ../protobuf && make

lintf:
	prototool lint ./protos
	prototool format -w
	go fmt ./...
	go mod tidy

genproto:
	prototool all

buildgo:
	go build ./...

pgfix:
ifeq ($(UNAME), Linux)
	for pb in  $(PG_FIX_CANDIDATES); \
	do \
	sed -ne 's/TableName/tableName/g' $$pb; \
	sed -ne 's/protobuf:"bytes,[0-9]*,opt,name=table_name,json=tableName,proto[0-9]*" json:"table_name,omitempty" pg:"table_name" //g' $$pb; \
	done
endif 
ifeq ($(UNAME), Darwin)
	for pb in  $(PG_FIX_CANDIDATES); \
	do \
	sed -i '' -e 's/TableName/tableName/g' $$pb; \
	sed -i '' -e 's/protobuf:"bytes,[0-9]*,opt,name=table_name,json=tableName,proto[0-9]*" json:"table_name,omitempty" pg:"table_name" //g' $$pb; \
	done
endif 

build: lintf genproto buildgo pgfix

test:
	dropdb --if-exists $(TEST_DB_NAME)
	createdb $(TEST_DB_NAME)
	go test -v ./... -args -db_url=$(TEST_DB_URL) -rd_url=$(TEST_RD_URL)
	dropdb $(TEST_DB_NAME)

