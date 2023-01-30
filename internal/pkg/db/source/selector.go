package source

import (
	"os"
	"log"
	"strings"

	"eventforward/internal/pkg/db/source/mongo"
)

func selector[T any](declname, name string) DB[T] {
	switch strings.ToLower(name) {
	case "mongo":
		instance, err := mongo.Setup[T]()
		if err != nil {
			log.Fatal(err)
		}
		return instance

	case "":
		log.Fatalf("Empty database source : %s\n", declname)

	default:
		log.Fatalf("Unknown database source : %s\n", name)
	}

	return nil
}

func DBFromName[T any](declname, name string) DB[T] {
	return selector[T](declname, name)
}

func DBFromEnv[T any]() DB[T] {
	name := os.Getenv("SOURCE_DATABASE")
	return selector[T]("env variable", name)
}
