PB = $(wildcard *.proto)
GO = $(PB:.proto=.pb.go)

all: $(GO)

%.pb.go: %.proto
		protoc --proto_path=$(GOPATH)/src:. --gogofaster_out=. $<

		src 目录下
		protoc --proto_path=. --gogofaster_out=. .\github.com\ipfs\go-bitswap\message\pb\message.proto
		protoc --proto_path=. --go_out=.\github.com\ipfs\go-bitswap\message\pb\. .\github.com\ipfs\go-bitswap\message\pb\message.proto

clean:
		rm -f *.pb.go
		rm -f *.go
