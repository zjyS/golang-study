package repo

import "example.com/web-service-gin/bean"

type AlbumLocalRepo struct {
	albums []bean.Album
}

func NewAlbumLocalRepo(albums *[]bean.Album) AlbumRepo {
	return &AlbumLocalRepo{*albums}
}

func (al *AlbumLocalRepo) GetAlbum(id string) (Album *bean.Album) {
	for _, a := range al.albums {
		if a.ID == id {
			Album = &a
			break
		}
	}
	return
}

func (al *AlbumLocalRepo) PutAlbum(a bean.Album) bool {
	isPut := true
	for _, oldA := range al.albums {
		if oldA.ID == a.ID {
			isPut = false
			break
		}
	}
	if isPut {
		al.albums = append(al.albums, a)
	}
	return isPut
}

func (al *AlbumLocalRepo) DeleteAlbum(id string) bool {
	newAlbums := make([]bean.Album, 0)

	exists := false

	for _, a := range al.albums {
		if a.ID != id {
			newAlbums = append(newAlbums, a)
			continue
		}
		exists = true
	}

	if exists {
		al.albums = newAlbums
	}

	return exists
}

func (al *AlbumLocalRepo) UpdateAlbum(a bean.Album) {
	isUpdate := false
	for _, oldA := range al.albums {
		if oldA.ID == a.ID {
			isUpdate = true
			break
		}
	}
	if isUpdate {
		al.DeleteAlbum(a.ID)
		al.PutAlbum(a)
	}
}

func (al *AlbumLocalRepo) GetAllEditable() *[]bean.Album {
	return &al.albums
}

func (al *AlbumLocalRepo) GetAll() []bean.Album {
	return al.albums
}

func (al *AlbumLocalRepo) Close() {}
