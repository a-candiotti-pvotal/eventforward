package destination

import (
	"log"
	"os"
	"strings"

	"eventforward/internal/pkg/db/destination/eventstore"
)

func selector[T any](declname, name string) DB[T] {
	switch strings.ToLower(name) {
	case "eventstore":
		instance, err := eventstore.Setup[T]()
		if err != nil {
			log.Fatal(err)
		}
		return instance

	case "":
		log.Fatalf("Empty database destination : %s\n", declname)

	default:
		log.Fatalf("Unknown database destination : %s\n", name)
	}

	return nil
}

func DBFromName[T any](declname, name string) DB[T] {
	return selector[T](declname, name)
}

func DBFromEnv[T any]() DB[T] {
	name := os.Getenv("DESTINATION_DATABASE")
	return selector[T]("env variable", name)
}
