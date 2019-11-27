FROM golang:1.13-stretch

MAINTAINER TRON-US <support@tron.network>
# Dockerfile.unit_testing will build an image to run the go unit tests.
# Use the regular Dockerfile to run a btfs daemon instead

#ENV PATH="/go/bin:${PATH}"

ENV PROTOC_VERSION=3.10.0
ENV GOLANG_PROTOBUF_VERSION=1.3.2
ENV PROTOTOOL_VERSION=1.9.0


# Install patch
RUN apt-get update && apt-get install -y patch

RUN apt-get install -y unzip

# install standard c++ implementation of protocol buffers
RUN wget --quiet https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip
RUN unzip protoc-${PROTOC_VERSION}-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include

# install golang proto package
RUN GO111MODULE=on go get \
  github.com/golang/protobuf/protoc-gen-go@v${GOLANG_PROTOBUF_VERSION} && \
  mv /go/bin/protoc-gen-go* /usr/local/bin/

# install prototool
RUN wget --quiet https://github.com/uber/prototool/releases/download/v${PROTOTOOL_VERSION}/prototool-Linux-x86_64.tar.gz
RUN tar -xf prototool-Linux-x86_64.tar.gz
RUN mv prototool/bin/prototool /usr/local/bin/prototool

ENV SRC_DIR /go/src/github.com/tron-us/go-btfs-common

# Download packages first so they can be cached.
COPY go.mod go.sum $SRC_DIR/
RUN cd $SRC_DIR \
  && go mod download

COPY . $SRC_DIR

#WORKDIR /go-btfs-common
WORKDIR $SRC_DIR

# by default lets run the go fmt, tidy and unit tests
CMD make trongogo build
