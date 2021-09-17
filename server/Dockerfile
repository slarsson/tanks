# build
FROM golang:1.16.0-alpine

RUN apk add --no-cache git
RUN apk add --no-cache make

WORKDIR /game

COPY . .

RUN go get -d -v ./...

RUN make amd64


# deploy
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=0 /game/build/server_amd64 .

EXPOSE 1337

CMD ["./server_amd64"]