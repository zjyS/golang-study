package repo

import "example.com/web-service-gin/bean"

type AlbumRepo interface {
	AlbumReadRepo
	AlbumWriteRepo
	Closer
}

type AlbumReadRepo interface {
	GetAlbum(id string) *bean.Album
	GetAll() []bean.Album
	GetAllEditable() *[]bean.Album
}

type AlbumWriteRepo interface {
	PutAlbum(a bean.Album) bool
	DeleteAlbum(id string) bool
	UpdateAlbum(a bean.Album)
}

type Closer interface {
	Close()
}
