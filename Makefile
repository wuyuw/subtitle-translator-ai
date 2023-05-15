# Build all by default, even if it's not first
.DEFAULT_GOAL := build

BUILD_NAME=sta

.PHONY: build
build:
	go build -o ${BUILD_NAME} .



	