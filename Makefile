all: build

build: deps install proto service fmt

deps: protoc-gen-go grpc-gateway glide

install:
	glide install

run:
	./book_parser

service:
	go build

fmt:
	go fmt ./

proto:
	protoc -I/usr/local/include -I. \
	-I$$GOPATH/src \
	-I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:. server/proto/*.proto
	protoc -I/usr/local/include -I. \
  -I$$GOPATH/src \
  -I$$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --grpc-gateway_out=logtostderr=true:.  server/proto/*.proto



protoc-gen-go:
	go get -u github.com/golang/protobuf/protoc-gen-go

grpc-gateway:
	go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

glide:
	go get -u github.com/Masterminds/glide
