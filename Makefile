default: lintf

SRC_DIR=./info ./ledger

install:
	brew install protobuf
	brew install prototool

lintf:
	prototool lint ./protos
	prototool format -w

build:
	prototool all
	go mod tidy
	for dir in $(SRC_DIR); \
	do \
	eval go build $$dir; \
	done
