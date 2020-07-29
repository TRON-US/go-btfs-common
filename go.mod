module github.com/tron-us/go-btfs-common

go 1.13

require (
	github.com/btcsuite/btcutil v0.0.0-20190425235716-9e5f4b9a998d
	github.com/ethereum/go-ethereum v1.9.17
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.3
	github.com/libp2p/go-libp2p-core v0.0.6
	github.com/stretchr/testify v1.4.0
	github.com/tron-us/go-common/v2 v2.1.1
	github.com/tron-us/protobuf v1.3.4
	go.uber.org/zap v1.14.1
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4 // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
	google.golang.org/genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc v1.25.1
)

replace github.com/libp2p/go-libp2p-core => github.com/TRON-US/go-libp2p-core v0.4.1
