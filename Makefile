default: lint

install:
	brew install protobuf
	brew install prototool

lint:
	prototool lint ./protos
