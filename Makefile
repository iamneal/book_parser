all: build

build: deps install proto service fmt

deps: protoc-gen-go glide

install:
	glide install

run:
	./book_parser

service:
	go build

fmt:
	go fmt ./

proto:
	cd server/proto && protoc --go_out=plugins=grpc:. *.proto && cd ../../

protoc-gen-go:
	go get -u github.com/golang/protobuf/protoc-gen-go

glide:
	go get -u github.com/Masterminds/glide
