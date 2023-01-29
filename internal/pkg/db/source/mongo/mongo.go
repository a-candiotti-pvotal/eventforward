package mongo

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"eventforward/internal/pkg/models"
)

type MongoDB struct {
	client *mongo.Client
}

func Setup() (*MongoDB, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("MONGO_URI env variable was empty")
	}

	// FIXME : singleton
	client, err := connect(mongoURI)
	if err != nil {
		return nil, err
	}
	m := &MongoDB{
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

func (m *MongoDB) ReadOperations(done chan struct{}, opChan chan<- *models.ChangeEvent, errChan chan<- error, namespace string) {
	ctx, cancel := context.WithCancel(context.Background())
	go func () {
		select {
		case <-done:
			cancel()
			return
		}
	}()

	snamespace := strings.Split(namespace, ".")
	if len(snamespace) != 2 {
		log.Fatalf("Malformated namespace, should be collection.name : %s\n", namespace)
	}

	targetDatabaseName := snamespace[0]
	targetCollectionName := snamespace[1]

	c := m.client.Database(targetDatabaseName).Collection(targetCollectionName)

	filter := bson.D{}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	var results []map[string]interface{}
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		cursor.Decode(&result)
		if err != nil {
			// FIXME : stop reading?
			log.Errorf("Decode fail on namespace %s.%s : %s", targetDatabaseName, targetCollectionName, err)
		} else {
			opChan <- &models.ChangeEvent{
				// FIXME : what about other fields?
				// OperationType: "insert", ? make sense?
				FullDocument: result,
			}
		}
	}
}
