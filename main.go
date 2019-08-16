package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var db *bbolt.DB
var channels []string

func main() {
	channels = []string{"jd", "taobao", "tmall", "suning"}

	var err error
	db, err = bbolt.Open("./db", 0666, nil)
	handle(err)
	defer db.Close()
	err = db.Batch(func(tx *bbolt.Tx) error {
		for _, v := range channels {
			_, err = tx.CreateBucketIfNotExists([]byte(v))
		}
		return nil
	})
	handle(err)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST("/:channel/:id", dealUpload)
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Review Collect Server(0.0.1). https://git.ixarea.com/ix64/review_collect")
	})
	r.GET("/:channel", dealDownload)
	err = r.Run("127.0.0.1:6481")
	handle(err)
}

func dealUpload(c *gin.Context) {
	channel := c.Param("channel")
	if !checkExist(channel, channels) {
		c.JSON(http.StatusNotAcceptable, gin.H{"status": false})
		return
	}
	id := c.Param("id")
	f, err := os.OpenFile("./"+channel+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		handle(err)
	} else {
		defer f.Close()
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	handle(err)

	line := bytes.Join([][]byte{[]byte(id), body, {'\n'}}, []byte(""))
	_, err = f.Write(line)
	handle(err)

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(channel))
		return b.Put([]byte(id), body)

	})
	handle(err)
	c.JSON(http.StatusOK, gin.H{"status": true})
}
func dealDownload(c *gin.Context) {
	channel := c.Param("channel")
	if !checkExist(channel, channels) {
		c.JSON(http.StatusNotAcceptable, gin.H{"status": false})
		return
	}
	var data []byte
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(channel))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			data = bytes.Join([][]byte{data, v}, []byte{'\n'})
		}
		return nil
	})
	handle(err)
	c.Data(http.StatusOK, "text/plain; charset=utf-8", data)
}
func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkExist(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
