package source

import (
	"log"

	"eventforward/internal/pkg/db/source/mongo"
)

// TODO : check name, check env variable?
func DBFromEnv[T any]() DB[T] {
	instance, err := mongo.Setup[T]()
	if err != nil {
		log.Fatal(err)
	}

	return instance
}
