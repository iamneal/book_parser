all: build

build: deps install proto-build service-build fmt

deps: protoc-gen-go glide

install:
	glide install

service-build:
	go build

fmt:
	go fmt ./

proto-build:
	cd server/proto && protoc --go_out=plugins=grpc:. *.proto && cd ../../

protoc-gen-go:
	go get -u github.com/golang/protobuf/protoc-gen-go

glide:
	go get -u github.com/Masterminds/glide
