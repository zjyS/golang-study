package client

import (
	"net"
	"net/url"
)
import "io"

type HttpClient struct {
	conn    net.Conn
	url     url.URL
	success func(io.Reader)
	failed  func(error)
}

func (httpClient *HttpClient) SendMessage(body []byte, success func(io.Reader), failed func(error)) error {
	httpClient.success = success
	httpClient.failed = failed
	_, err := httpClient.conn.Write(body)
	if err != nil {
		failed(err)
		return err
	}
	success(httpClient.conn)
	return nil
}
