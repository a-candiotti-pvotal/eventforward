package mongo

import (
	"context"
	"errors"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB[T any] struct {
	client *mongo.Client
}

func Setup[T any]() (*MongoDB[T], error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("MONGO_URI env variable was empty")
	}

	// FIXME : singleton
	client, err := connect(mongoURI)
	if err != nil {
		return nil, err
	}
	m := &MongoDB[T]{
		client: client,
	}
	return m, nil
}

func connect(uri string) (*mongo.Client, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)

	log.Println("Connecting to mongo")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Println("Connected successfully to mongo")
	return client, nil
}

func (m *MongoDB[T]) ReadOperations(done chan struct{}, opChan chan<- *T, errChan chan<- error, database, collection string) {
	c := m.client.Database(database).Collection(collection)

	ctx, cancel := context.WithCancel(context.Background())

	filter := bson.D{}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	go func () {
		select {
		case <-done:
			cancel()
			cursor.Close(context.Background())
			return
		}
	}()

	for _, result := range results {
		cursor.Decode(&result)
		if err != nil {
			// FIXME : stop reading? warning?
			log.Errorf("Decode fail on namespace %s.%s : %s", database, collection, err)
		} else {
			opChan <- &result
		}
	}
}
