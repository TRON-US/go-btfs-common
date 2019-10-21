default: lintf

PROTO_FILES=./protos/node/node.proto \
# 			./protos/status/status.proto \
# 			./protos/guard/guard.proto \
 			./protos/escrow/escrow.proto \

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
	eval protoc --go_out=plugins=grpc:. $$proto ; \
	done
