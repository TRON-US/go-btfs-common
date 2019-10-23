default: lintf

PROTO_FILES=./protos/node/node.proto \
			./protos/status/status.proto \
# 			./protos/guard/guard.proto \
# 			./protos/escrow/escrow.proto \

PB_OUT_PATH=$$GOPATH/src

install:
	brew install protobuf
	brew install prototool

lintf:
	prototool lint ./protos
	prototool format -w

build: lintf
# 	TODO: fix and use prototool all instead
	for proto in  $(PROTO_FILES); \
	do \
	eval protoc -I. --go_out=plugins=grpc:$(PB_OUT_PATH) $$proto ; \
	done