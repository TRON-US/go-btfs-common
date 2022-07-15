UNAME := $(shell uname)
default: lintf

TEST_DB_NAME ?= runtime
TEST_DB_STATUS ?= db_status
TEST_DB_GUARD ?= db_guard
TEST_DB_USER ?= `whoami`
TEST_DB_HOSTNAME ?= localhost
TEST_DB_PORT ?= 5432
TEST_DB_URL="postgresql://$(TEST_DB_USER)@$(TEST_DB_HOSTNAME):$(TEST_DB_PORT)/$(TEST_DB_NAME)"
TEST_DB_URL_STATUS="postgresql://$(TEST_DB_USER)@$(TEST_DB_HOSTNAME):$(TEST_DB_PORT)/$(TEST_DB_STATUS)"
TEST_DB_URL_GUARD="postgresql://$(TEST_DB_USER)@$(TEST_DB_HOSTNAME):$(TEST_DB_PORT)/$(TEST_DB_GUARD)"
TEST_RD_NAME ?= 4
TEST_RD_USER ?= `whoami`
TEST_RD_HOSTNAME ?= localhost
TEST_RD_PORT ?= 6379
TEST_RD_URL="redis://$(TEST_RD_USER):@$(TEST_RD_HOSTNAME):$(TEST_RD_PORT)/$(TEST_RD_NAME)"
DOCKER_TEST_RD_URL="redis://$(TEST_RD_HOSTNAME):$(TEST_RD_PORT)"

PG_FIX_CANDIDATES=./protos/node/node.pb.go \
			./protos/status/status.pb.go \
			./protos/escrow/escrow.pb.go \
			./protos/guard/guard.pb.go \
			./protos/online/online.pb.go \

install: brew trongogo

brew:
	brew install go protobuf prototool postgresql redis

trongogo:
	cd ../ && git clone https://github.com/googleapis/googleapis.git && cd go-btfs-common || true
	cd ../ && git clone https://github.com/TRON-US/protobuf || true
	cd ../protobuf && git checkout master && git pull && make
	cd ../googleapis && git checkout master && git pull

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
	sed -in 's/TableName/tableName/g' $$pb; \
	sed -in 's/protobuf:"bytes,[0-9]*,opt,name=table_name,json=tableName,proto[0-9]*" json:"table_name,omitempty" pg:"table_name" //g' $$pb; \
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
	brew services start postgresql
	brew services start redis
	dropdb --if-exists $(TEST_DB_NAME)
	dropdb --if-exists $(TEST_DB_GUARD)
	dropdb --if-exists $(TEST_DB_STATUS)
	createdb $(TEST_DB_NAME)
	createdb $(TEST_DB_GUARD)
	createdb $(TEST_DB_STATUS)
	TEST_DB_URL=$(TEST_DB_URL) TEST_DB_URL_STATUS=$(TEST_DB_URL_STATUS) TEST_DB_URL_GUARD=$(TEST_DB_URL_GUARD) TEST_RD_URL=$(TEST_RD_URL) go test -v ./...
	dropdb $(TEST_DB_NAME)
	dropdb $(TEST_DB_STATUS)
	dropdb $(TEST_DB_GUARD)
	brew services stop postgresql
	brew services stop redis

test_docker:
	sleep 10
	TEST_DB_URL=$(TEST_DB_URL) TEST_DB_URL_STATUS=$(TEST_DB_URL_STATUS) TEST_DB_URL_GUARD=$(TEST_DB_URL_GUARD) TEST_RD_URL=$(TEST_RD_URL) go test -v ./...

test_git_diff_protos:
	bin/test-git-diff-protos
.PHONY: test_git_dif_protos
