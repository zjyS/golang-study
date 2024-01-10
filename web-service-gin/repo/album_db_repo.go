package repo

import (
	"database/sql"
	"fmt"

	"example.com/web-service-gin/bean"
	_ "github.com/go-sql-driver/mysql"
)

var dbRepo *AlbumDBRepo = nil

type AlbumDBRepo struct {
	db    *sql.DB
	table string
}

func NewAlbumDBRepo(driver string, dataSource string) AlbumRepo {
	if dbRepo == nil {
		dbRepo = &AlbumDBRepo{}
		dbRepo.init(driver, dataSource)
		dbRepo.table = "album"
	}
	return dbRepo
}

func (al *AlbumDBRepo) init(driver string, dataSource string) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		panic(err.Error())
	}
	al.db = db
}

func (al *AlbumDBRepo) GetAlbum(id string) *bean.Album {
	sqlStr := fmt.Sprintf("SELECT * FROM `%s` WHERE `id` = ?", al.table)
	stmt, err := al.db.Prepare(sqlStr)
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	// 执行查询
	rows, err := stmt.Query(id) // 使用参数值替换替换项
	if err != nil {
		panic(err.Error())
	}
	defer func(rs *sql.Rows) {
		err := rs.Close()
		if err != nil {

		}
	}(rows)

	var album *bean.Album
	for rows.Next() {
		var (
			ID, Title, Artist string
		)
		var Price float64
		err := rows.Scan(&ID, &Title, &Artist, &Price)
		if err != nil {
			panic(err.Error())
		}
		album = &bean.Album{ID: ID, Title: Title, Artist: Artist, Price: Price}
		break
	}
	if err = rows.Err(); err != nil {
		panic(err.Error())
	}
	return album
}

func (al *AlbumDBRepo) PutAlbum(a bean.Album) bool {
	album := al.GetAlbum(a.ID)
	if album != nil {
		return false
	}
	sqlStr := fmt.Sprintf("INSERT INTO %s VALUES(?,?,?,?)", al.table)
	stmt, err := al.db.Prepare(sqlStr)
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)
	ret, err := stmt.Exec(a.ID, a.Title, a.Artist, a.Price)
	if err != nil {
		panic(err.Error())
	}

	es, err := ret.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	return es > 0
}

func (al *AlbumDBRepo) DeleteAlbum(id string) bool {
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE `id` = ?", al.table)
	stmt, err := al.db.Prepare(sqlStr)
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	ret, err := stmt.Exec(id)
	if err != nil {
		panic(err.Error())
	}

	es, err := ret.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	return es > 0
}

func (al *AlbumDBRepo) UpdateAlbum(a bean.Album) {
	sqlStr := fmt.Sprintf("UPDATE %s SET `title` = ?, `artist` = ?, `price` = ? WHERE `id` = ?", al.table)
	stmt, err := al.db.Prepare(sqlStr)
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	ret, err := stmt.Exec(a.Title, a.Artist, a.Price, a.ID)
	if err != nil {
		panic(err.Error())
	}

	_, err = ret.RowsAffected()
	if err != nil {
		panic(err.Error())
	}
}

func (al *AlbumDBRepo) GetAll() []bean.Album {
	sqlStr := fmt.Sprintf("SELECT * FROM %s", al.table)
	rows, err := al.db.Query(sqlStr)
	if err != nil {
		panic(err.Error())
	}
	defer func(rs *sql.Rows) {
		err := rs.Close()
		if err != nil {

		}
	}(rows)

	var albums = make([]bean.Album, 0)
	for rows.Next() {
		var (
			ID, Title, Artist string
		)
		var Price float64
		err := rows.Scan(&ID, &Title, &Artist, &Price)
		if err != nil {
			panic(err.Error())
		}
		albums = append(albums, bean.Album{ID: ID, Title: Title, Artist: Artist, Price: Price})
	}
	return albums
}

func (al *AlbumDBRepo) GetAllEditable() *[]bean.Album {
	sqlStr := fmt.Sprintf("SELECT * FROM %s", al.table)
	rows, err := al.db.Query(sqlStr)
	if err != nil {
		panic(err.Error())
	}
	defer func(rs *sql.Rows) {
		err := rs.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)

	var albums = make([]bean.Album, 0)
	for rows.Next() {
		var (
			ID, Title, Artist string
		)
		var Price float64
		err := rows.Scan(&ID, &Title, &Artist, &Price)
		if err != nil {
			panic(err.Error())
		}
		albums = append(albums, bean.Album{ID: ID, Title: Title, Artist: Artist, Price: Price})
	}
	return &albums
}

func (al *AlbumDBRepo) Close() {
	err := al.db.Close()
	if err != nil {
		return
	}
	fmt.Printf("db connect closed")
}
