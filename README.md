# go-btfs-common
Common packages among various go-btfs projects

## File Structure

```
.
├── LICENSE
├── Makefile
├── README.md
├── go.mod
├── info
│   └── node.go (deprecating)
├── protos
│   └── node
│       ├── node.pb.go
│       └── node.proto
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
