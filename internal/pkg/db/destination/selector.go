package destination

import (
	"log"

	"eventforward/internal/pkg/db/destination/eventstore"
)

// TODO : switch on name in env variable?
func DBFromEnv[T any]() DB[T] {
	instance, err := eventstore.Setup[T]()
	if err != nil {
		log.Fatal(err)
	}

	return instance
}
