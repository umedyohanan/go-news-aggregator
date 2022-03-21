package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"newsaggregator/pkg/storage"
)

type Store struct {
	db *mongo.Client
}

const (
	databaseName   = "posts"
	collectionName = "posts"
)

// New Constructor, gets connection string.
func New(constr string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	db, err := mongo.Connect(context.Background(), mongoOpts)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

// Posts Return list of posts from database.
// Gets number of posts to return.
func (s *Store) Posts(n int64) ([]storage.Post, error) {
	collection := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	findOpt := options.Find()
	findOpt.SetSort(bson.D{{"_id", 1}})
	findOpt.SetLimit(n)
	cur, err := collection.Find(context.Background(), filter, findOpt)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	var data []storage.Post
	for cur.Next(context.Background()) {
		var p storage.Post
		err := cur.Decode(&p)
		if err != nil {
			return nil, err
		}
		data = append(data, p)
	}
	return data, cur.Err()
}

// SavePosts Saves given list of posts using mongo's InsertMany
func (s *Store) SavePosts(posts []storage.Post) error {
	var docs []interface{}
	for _, t := range posts {
		docs = append(docs, t)
	}
	collection := s.db.Database(databaseName).Collection(collectionName)
	_, err := collection.InsertMany(context.Background(), docs)
	return err
}
