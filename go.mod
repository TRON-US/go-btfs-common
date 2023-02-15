module github.com/tron-us/go-btfs-common

go 1.13

require (
	github.com/btcsuite/btcutil v1.0.2
	github.com/ethereum/go-ethereum v1.9.24
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hypnoglow/go-pg-monitor v0.1.0
	github.com/hypnoglow/go-pg-monitor/gopgv9 v0.1.0
	github.com/libp2p/go-libp2p-core v0.0.6
	github.com/prometheus/client_golang v1.11.1
	github.com/stretchr/testify v1.6.1
	github.com/tron-us/go-common/v2 v2.3.0
	github.com/tron-us/protobuf v1.3.4
	go.uber.org/zap v1.16.0
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.34.0
)

replace github.com/libp2p/go-libp2p-core => github.com/TRON-US/go-libp2p-core v0.4.1
