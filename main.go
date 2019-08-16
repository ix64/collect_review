package main

import (
	"bytes"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
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
	_ = godotenv.Load(".env")
	log.Println("Write Upload Log At (LOG_PATH): ", os.Getenv("LOG_PATH"))

	var err error
	log.Println("Use Database (DB_PATH): ", os.Getenv("DB_PATH"))
	db, err = bbolt.Open(os.Getenv("DB_PATH"), 0666, nil)
	handle(err)
	defer db.Close()
	err = db.Batch(func(tx *bbolt.Tx) error {
		for _, v := range channels {
			_, err = tx.CreateBucketIfNotExists([]byte(v))
		}
		return nil
	})
	handle(err)

	r := httprouter.New()
	r.POST("/:channel/:id", dealUpload)
	r.GET("/", dealHome)
	r.GET("/:channel", dealDownload)
	log.Print("Listen at (BIND): ", os.Getenv("BIND"))
	err = http.ListenAndServe(os.Getenv("BIND"), r)
	handle(err)
}

func dealUpload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	channel := ps.ByName("channel")
	if !checkExist(channel, channels) {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	id := ps.ByName("id")
	f, err := os.OpenFile(os.Getenv("LOG_PATH")+channel+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		handle(err)
	} else {
		defer f.Close()
	}
	body, err := ioutil.ReadAll(r.Body)
	handle(err)

	line := bytes.Join([][]byte{[]byte(id), body, {'\n'}}, []byte(""))
	_, err = f.Write(line)
	handle(err)

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(channel))
		return b.Put([]byte(id), body)

	})
	handle(err)
	w.WriteHeader(http.StatusOK)
}
func dealDownload(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	channel := ps.ByName("channel")
	if !checkExist(channel, channels) {
		w.WriteHeader(http.StatusNotAcceptable)
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
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(data)
}
func dealHome(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	data := "Review Collect Server(0.0.1). https://githu.com/ix64/review_collect"
	_, _ = w.Write([]byte(data))
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
