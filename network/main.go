package main

import (
	"bufio"
	"encoding/json"
	"example.com/web-service-gin/bean"
	"fmt"
	"io"
	"net"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")
	getById(conn)
}

func getById(conn net.Conn) {
	_, err := fmt.Fprintf(conn, "GET /albums/1 HTTP/1.0\r\n\r\n")
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	reader := bufio.NewReader(conn)
	var body = make([]byte, reader.Size())
	for i := 0; ; i = i + 1 {
		line, _, err := reader.ReadLine()
		if err == io.EOF || err != nil || line == nil {
			break
		}
		if i < 5 {
			continue
		}
		body = append(body, line...)
	}
	fmt.Printf("%s\n", body)
}

func post(conn net.Conn) {
	a := &bean.Album{ID: "7", Title: "GoLand", Artist: "John", Price: 88.99}
	marshal, err := json.Marshal(a)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(conn, "POST /albums HTTP/1.0\r\n")
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(conn, "Host: localhost:8080\r\n")
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(conn, "Content-Type: application/json; charset=utf-8\r\n")
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(conn, "Content-Length: %d\r\n\r\n", len(marshal))
	if err != nil {
		return
	}

	size, err := conn.Write(marshal)
	if err != nil {
		return
	}

	fmt.Printf("send data size : %d\n", size)
}
