package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"newsaggregator/pkg/api"
	"newsaggregator/pkg/rss"
	"newsaggregator/pkg/storage"
	"newsaggregator/pkg/storage/mongo"
	"time"
)

type server struct {
	db  storage.Interface
	api *api.API
}

// app config
type config struct {
	Urls          []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

func main() {
	var srv server
	var conf config
	// reading and decoding config file
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(configFile, &conf)
	if err != nil {
		log.Fatal(err)
	}
	// initializing database
	db, err := mongo.New("mongodb://localhost:27017/")
	if err != nil {
		log.Fatal(err)
	}
	// parsing news feeds in separate goroutine for each link
	chPosts := make(chan []storage.Post)
	chErrs := make(chan error)
	for _, url := range conf.Urls {
		go parseFeed(url, chPosts, chErrs, conf.RequestPeriod)
	}
	// saving posts from feed to db
	go func() {
		for posts := range chPosts {
			db.SavePosts(posts)
		}
	}()

	go func() {
		for err := range chErrs {
			log.Println("error:", err)
		}
	}()
	// starting webapp
	srv.db = db
	srv.api = api.New(srv.db)
	err = http.ListenAndServe(":80", srv.api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// parseFeed reads feed asynchronously.
// Decoded news and error are written to channels.
func parseFeed(url string, posts chan<- []storage.Post, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
