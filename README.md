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
├── config
│   └── common.go
├── crypto
│   └── crypto.go
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
├── prototool.yaml
├── tree.txt
└── utils
    ├── grpc_health_check_provider.go
    ├── runtime.go
    └── runtime_test.go
```

## Install Tools

#### install TRON version(pg support) protobuf

```
cd github.com/tron-us
```
```
git clone https://github.com/TRON-US/protobuf
```
```
cd protobuf && make
```

#### other tools

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

## Use docker container to run 'git diff protos/'

```
$ docker build -f Dockerfile -t "go-btfs-common" .
$ docker run -i go-btfs-common
```

## Run interactive bash inside docker container for diagnosis

```
$ docker build -f Dockerfile -t "go-btfs-common" .
$ docker run -it go-btfs-common /bin/bash
```

