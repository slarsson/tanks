# game

## Client
TODO

## Server
TODO

### Makefile:

Build server + wasm
```
make build
```

Build wasm part only
```
make wasm
```

Run server
```
make run
```
## Misc



Go target = wasm
```
GOOS=js GOARCH=wasm go build -o main.wasm
```

https://github.com/golang/go/wiki/GoArm

https://developer.valvesoftware.com/wiki/Source_Multiplayer_Networking