.PHONY: build wasm all

all: build

build:
	go build -o build/server .
	GOOS=js GOARCH=wasm go build -o build/main.wasm ./wasm
	cp ./build/main.wasm ../client/assets

wasm:
	GOOS=js GOARCH=wasm go build -o build/main.wasm ./wasm
	cp ./build/main.wasm ../client/assets

arm:
	GOOS=linux GOARCH=arm go build -o build/server_arm32 .

amd64:
	GOOS=linux GOARCH=amd64 go build -o build/server_amd64 .

run:
	go run .