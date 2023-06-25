package main

import (
	// _ "github.com/envoyproxy/protoc-gen-validate"
	"github.com/livekit/protocol/livekit"
	_ "github.com/livekit/protocol/livekit"
	"google.golang.org/protobuf/proto"
)

func main() {
	marshal, err := proto.Marshal(nil)
	if err != nil {
		println("%v", marshal, err, livekit.WebhookEvent{})
	}
	println("Hello World")
}
