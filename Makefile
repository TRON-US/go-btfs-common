default: lintf

PROTO_FILES=./protos/node/node.proto \

 			./protos/guard/guard.proto \
 			./protos/shared/serverstatus.proto \
			./protos/escrow/escrow.proto \
			./protos/ledger/ledger.proto \
# 			./protos/status/status.proto \
# 			./protos/guard/guard.proto \



install:
	brew install protobuf
	brew install prototool

lintf:
	prototool lint ./protos
	prototool format -w

build: lintf
	for proto in  $(PROTO_FILES); \
	do \
	eval protoc -I. --go_out=plugins=grpc:. $$proto ; \
	done
