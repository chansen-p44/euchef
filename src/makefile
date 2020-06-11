export GOOS=linux
export GOARCH=amd64

POST_OUTPUT = ../bin/post_lambda
DELETE_OUTPUT = ../bin/delete_lambda

build: makefile
	@go build -ldflags="-s -w" -o $(POST_OUTPUT) ./post_lambda
	@go build -ldflags="-s -w" -o $(DELETE_OUTPUT) ./delete_lambda

.PHONY: build

