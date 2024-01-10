package protocol

import (
	"example.com/web-service-gin/bean"
	"google.golang.org/protobuf/proto"
	"os"
)

func WriteMessage(album *Album) {
	al := bean.Album{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99}
	album.Id = &al.ID
	album.Title = &al.Title
	album.Artist = &al.Artist
	album.Price = &al.Price
	marshal, err := proto.Marshal(album)
	if err != nil {
		return
	}

	err = os.WriteFile("album", marshal, 0644)
	if err != nil {
		return
	}
}

func ReadMessage() *Album {
	file, err := os.ReadFile("album")
	if err != nil {
		return nil
	}
	album := &Album{}
	err = proto.Unmarshal(file, album)
	if err != nil {
		return nil
	}
	return album
}
