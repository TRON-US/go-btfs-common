# go-btfs-common

Common packages among various go-btfs projects.

| golang | protobuf  | prototool |
|--------|-----------|-----------|
| 1.13   | 3.10.0    | 1.9.0     |

## File Structure

```
.
├── LICENSE
├── Makefile
├── README.md
├── go.mod
├── go.sum
├── info
│   └── node.go (deprecating)
├── ledger
│   └── ledger.go
├── protos
│   ├── escrow
│   │   ├── escrow.pb.go
│   │   └── escrow.proto
│   ├── guard
│   │   ├── guard.pb.go
│   │   └── guard.proto
│   ├── ledger
│   │   ├── ledger.pb.go
│   │   └── ledger.proto
│   ├── node
│   │   ├── node.pb.go
│   │   └── node.proto
│   ├── shared
│   │   ├── shared.pb.go
│   │   └── shared.proto
│   └── status
│       ├── status.pb.go
│       └── status.proto
└── prototool.yaml
```

## Install Tools

```
make install
```

## Lint and Format

```
make
```
or
```
make lintf
```

## Build/Complie

```
make build
```
