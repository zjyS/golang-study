package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"example.com/web-service-gin/bean"
	"example.com/web-service-gin/repo"
	"github.com/gin-gonic/gin"
)

// albums slice to seed record Album data.
var albums = []bean.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

var albumLocalHolder = repo.NewAlbumLocalRepo(&albums)

const (
	driver      = "mysql"
	data_source = "test:Idle123$@tcp(localhost:3306)/mydb"
)

var albumDBHolder = repo.NewAlbumDBRepo(driver, data_source)

func getAlbumHolder() repo.AlbumRepo {
	return albumDBHolder
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, getAlbumHolder().GetAll())
}

// postAlbums adds an Album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum bean.Album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new Album to the slice.
	ret := getAlbumHolder().PutAlbum(newAlbum)
	if ret {
		c.IndentedJSON(http.StatusCreated, newAlbum)
		return
	}

	c.IndentedJSON(http.StatusBadRequest, "Duplicate data!")
}

// getAlbumByID locates the Album whose ID value matches the id
// parameter sent by the client, then returns that Album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an Album whose ID value matches the parameter.
	Album := getAlbumHolder().GetAlbum(id)
	if Album != nil {
		c.IndentedJSON(http.StatusOK, Album)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
}

func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")

	ret := getAlbumHolder().DeleteAlbum(id)

	if !ret {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, getAlbumHolder().GetAll())
}

func deleteAllAlbum(c *gin.Context) {
	p := getAlbumHolder().GetAllEditable()
	p2 := make([]bean.Album, 0, 3)
	*p = p2
	c.IndentedJSON(http.StatusOK, getAlbumHolder().GetAll())
}

func init() {
	fmt.Println("init finished!")
}

func main() {
	// 创建一个通道用于接收信号
	sigChan := make(chan os.Signal, 1)

	// 将指定的信号发送到通道
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 监听信号
	go func() {
		// 接收到信号后执行的逻辑
		sig := <-sigChan
		fmt.Printf("接收到信号：%v\n", sig)

		// 在这里执行进程销毁前的清理操作
		getAlbumHolder().Close()
		// 结束程序
		os.Exit(0)
	}()

	// 主程序继续执行其他操作
	// ...
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.DELETE("delete/:id", deleteAlbumByID)
	router.DELETE("delete", deleteAllAlbum)

	err := router.Run("localhost:8080")
	if err != nil {
		return
	}
}
