package mongo

import (
	"strings"
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

	"eventforward/internal/pkg/models"
)

// TODO : do a common function with ReadOperations that have different behavior for watch and read?
// something like that, code is very close

func (m *MongoDB) WatchOperations(done chan struct{}, opChan chan<- *models.ChangeEvent, _ chan<- error, namespace string) {
	snamespace := strings.Split(namespace, ".")
	if len(snamespace) != 2 {
		log.Fatalf("Malformated namespace, should be collection.name : %s\n", namespace)
	}

	targetDatabaseName := snamespace[0]
	targetCollectionName := snamespace[1]

	c := m.client.Database(targetDatabaseName).Collection(targetCollectionName)

	ctx, cancel := context.WithCancel(context.Background())
	changeStream, err := c.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		log.Fatal(err)
	}

	go func () {
		select {
		case <-done:
			changeStream.Close(ctx)
			cancel()
			return
		}
	}()

	for changeStream.Next(ctx) {
		changeEvent := &models.ChangeEvent{}
		err := changeStream.Decode(changeEvent)
		if err != nil {
			log.Errorf("Decode fail on namespace %s.%s : %s", targetDatabaseName, targetCollectionName, err)
		} else {
			opChan <- changeEvent
		}
	}
}
