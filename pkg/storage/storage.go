package storage

type Post struct {
	Title   string
	Content string
	PubTime int64
	Link    string
}

type Interface interface {
	Posts(n int64) ([]Post, error) // get posts
	SavePosts([]Post) error        // adding posts from rss
}
