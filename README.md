# go-btfs-common
Common packages among various go-btfs projects

## File Structure

```
.
├── LICENSE
├── Makefile
├── README.md
├── go.mod
├── info (deprecating)
│   └── node.go
├── protos
│   ├── node
│   │   ├── node.pb.go
│   │   └── node.proto
│   └── status
│       ├── status.pb.go
│       └── status.proto
├── prototool.yaml
└── tree.txt
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
