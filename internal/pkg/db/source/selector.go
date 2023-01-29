package source

import (
	"log"

	"eventforward/internal/pkg/db/source/mongo"
)

// TODO : check name, check env variable?
func DBFromEnv() DB {
	instance, err := mongo.Setup()
	if err != nil {
		log.Fatal(err)
	}

	return instance
}
