package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/web-service-gin/bean"
)

func main() {
	GetByID()
}

func GetAll() {
	get, err := http.Get("http://localhost:8080/albums")
	if err != nil {
		return
	}
	body := handleResp(get)
	a := make([]bean.Album, 0)
	err = json.Unmarshal(body, &a)
	if err != nil {
		return
	}
	fmt.Println(a)

	marshal, err := json.Marshal(a)
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("%s", marshal))
}

func GetByID() {
	get, err := http.Get("http://localhost:8080/albums/1")
	if err != nil {
		return
	}
	body := handleResp(get)
	a := bean.Album{}
	err = json.Unmarshal(body, &a)
	if err != nil {
		return
	}
	fmt.Println(a)

	marshal, err := json.Marshal(a)
	if err != nil {
		return
	}
	fmt.Println(fmt.Sprintf("%s", marshal))
}

func Post() {
	a := bean.Album{ID: "5", Title: "xxx", Artist: "xxx", Price: 67.99}
	marshal, err := json.Marshal(a)
	if err != nil {
		return
	}
	buffer := bytes.Buffer{}
	buffer.Write(marshal)
	post, err := http.Post("http://localhost:8080/albums", "application/json", &buffer)
	if err != nil {
		return
	}
	body := handleResp(post)
	fmt.Println(fmt.Sprintf("%s", body))
}

func handleResp(resp *http.Response) []byte {
	bodySize := int64(0)
	body := make([]byte, 0, resp.ContentLength)
	bufferSize := 256
	for buffer := make([]byte, bufferSize); bodySize < resp.ContentLength; {
		num, err := resp.Body.Read(buffer)
		if num == 0 || (err != nil && err != io.EOF) {
			return body
		}
		body = append(body, buffer[:num]...)
		bodySize = bodySize + int64(num)
		if err == io.EOF {
			break
		}
	}
	return body
}
