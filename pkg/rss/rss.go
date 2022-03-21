package rss

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"newsaggregator/pkg/storage"
	"strings"
	"time"
)

type Rss struct {
	Chanel Channel `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubdate"`
}

// Rss feed parser, returns posts as slice of storage.Post.
func Parse(url string) ([]storage.Post, error) {
	var posts []storage.Post
	resp, err := http.Get("https://habr.com/ru/rss/hub/go/all/?fl=ru")
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	rss := Rss{}

	err = xml.Unmarshal(data, &rss)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range rss.Chanel.Items {
		var post storage.Post

		post.Title = item.Title
		post.Content = item.Description
		post.Link = item.Link

		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			post.PubTime = t.Unix()
		}
		posts = append(posts, post)
	}

	return posts, nil
}
