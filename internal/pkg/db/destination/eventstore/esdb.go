package eventstore

import (
	"time"
	"context"
	"encoding/json"
	"log"
	"os"
	"errors"

	"github.com/EventStore/EventStore-Client-Go/v3/esdb"
)

type EventStoreDB[T any] struct {
	client *esdb.Client
}

func Setup[T any]() (*EventStoreDB[T], error) {
	esdbURI := os.Getenv("ESDB_URI")
	if esdbURI == "" {
		return nil, errors.New("ESDB_URI env variable was empty")
	}

	client, err := connect(esdbURI)
	if err != nil {
		return nil, err
	}
	e := &EventStoreDB[T]{ client: client }
	return e, nil
}

func connect(uri string) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(uri)

	if err != nil {
		return nil, err
	}

	db, err := esdb.NewClient(settings)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// 	TODO : can we make that concurrent?
func marshalEvents[T Free](eventType string, events []T) ([]esdb.EventData, error) {
	results := []esdb.EventData{}

	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}

		// FIXME : dont send Log event itself, event.Doc instead?
		//       	 what about event.Update?
		eventData := esdb.EventData{
			//			EventID: , TODO : set me
			ContentType: esdb.ContentTypeJson,
			EventType:   eventType,
			Data:        data,
		}

		results = append(results, eventData)
	}

	return results, nil
}

// Hack with Generics to avoid transforming to interface{} at runtime
type Free interface{}

func appendBatchOfEventsToStream[T Free](e *esdb.Client, events []T, eventType, streamName string) error {
	eventDatas, err := marshalEvents(eventType, events)
	if err != nil {
		return err
	}

	_, err = e.AppendToStream(context.Background(), streamName, esdb.AppendToStreamOptions{}, eventDatas...)
	return err
}

const BufferSize = 500

func (e *EventStoreDB[T]) SendOperations(done chan struct{}, opChan <-chan *T, _ chan<- error, stream, eventType string) {
	total := 0
	nbr := 0
	last := time.Now()

	buffer := []*T{}

	for {
		select {
		case <-done:
			return

		case event := <-opChan:
			// TODO : retry? on error what to do? reconnect?
			// use https://github.com/cenkalti/backoff?
//			log.Printf("%s since last execution\n", time.Since(last))

			buffer = append(buffer, event)

			if len(buffer) >= BufferSize {
				// FIXME : how to get event type with generics?
				// event.Ns.Db
				err := appendBatchOfEventsToStream(e.client, buffer, eventType, stream)
				if err != nil {
					log.Println(err)
					//				errChan <- err
				}

				nbr += len(buffer)
				buffer = []*T{}
			}

			if time.Since(last) >= 1 * time.Second {
				total += nbr
				log.Printf("DEBUG : Sending %d/s on a total of %d\n", nbr, total)
				last = time.Now()
				nbr = 0
			}
			break
		}
	}
}
