package destination

import (
	"log"

	"eventforward/internal/pkg/db/destination/eventstore"
)

// TODO : switch on name in env variable?
func DBFromEnv() DB {
	instance, err := eventstore.Setup()
	if err != nil {
		log.Fatal(err)
	}

	return instance
}
