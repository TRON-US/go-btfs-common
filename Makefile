default: lintf

install:
	brew install protobuf
	brew install prototool

lintf:
	prototool lint ./protos
	prototool format -w

build:
	prototool all
	go mod tidy
	go build ./...
