package mongo

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO : do a common function with ReadOperations that have different behavior for watch and read?
// something like that, code is very close

func (m *MongoDB[T]) WatchOperations(done chan struct{}, opChan chan<- *T, _ chan<- error, database, collection string) {
	c := m.client.Database(database).Collection(collection)

	ctx, cancel := context.WithCancel(context.Background())

	changeStream, err := c.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		log.Fatal(err)
	}

	go func () {
		select {
		case <-done:
			cancel()
			changeStream.Close(context.Background())
			return
		}
	}()

	for changeStream.Next(ctx) {
		var result T
		err := changeStream.Decode(&result)
		if err != nil {
			log.Errorf("Decode fail on namespace %s.%s : %s", database, collection, err)
		} else {
			opChan <- &result
		}
	}
}
