package main

import (
	"example.com/protobuf/protocol"
	"fmt"
)

func main() {
	message := protocol.ReadMessage()
	fmt.Println(message)
}
