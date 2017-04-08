all: deps install build fmt

build: proto web-server web-app

deps: protoc-gen-go grpc-gateway glide yarn create-react-app

install:
	glide install
	cd ./server/web_app && yarn install

run:
	./book_parser

web-server:
	go build

web-app:
	cd ./server/web_app && npm run build

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

yarn:
	curl -o- -L https://yarnpkg.com/install.sh | bash

create-react-app:
	npm install -g create-react-app
