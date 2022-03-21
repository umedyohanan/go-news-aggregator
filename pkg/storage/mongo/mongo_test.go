package mongo

import (
	"math/rand"
	"newsaggregator/pkg/storage"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	_, err := New("mongodb://localhost:27017/")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_News(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	posts := []storage.Post{
		{
			Title: "Test Post",
			Link:  strconv.Itoa(rand.Intn(1_000_000_000)),
		},
	}
	db, err := New("mongodb://localhost:27017/")
	if err != nil {
		t.Fatal(err)
	}
	err = db.SavePosts(posts)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.Posts(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
}
