package main

/*
#cgo CFLAGS: -I./libs
#cgo LDFLAGS: -L. -lclient
#include "client.h"
*/
import "C"
import _ "fmt"

func main() {
	C.get()
}
