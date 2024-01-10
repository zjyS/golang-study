package client

import "io"

type Client interface {
	SendMessage(body []byte, success func(io.Reader), failed func(error)) error
}
