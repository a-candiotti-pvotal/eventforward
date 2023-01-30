package main

import (
	"log"

	publicmodels "eventforward/pkg/models"
	"eventforward/internal/eventforward"
	"eventforward/internal/pkg/models"
)

// FIXME : DEBUG
// FIXME : read from yaml/json
var decls = []publicmodels.ForwardDecl{
	{
		Name: "mongo-eventstore",
		From: publicmodels.ForwardDeclPoint{
			Type: "mongo",
			Database: "profiling",
			Table: "bet",
		},
		To: publicmodels.ForwardDeclPoint{
			Type: "eventstore",
			Database: "profiling",
			Table: "bet",
		},
		Watch: true,
	},
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// FIXME : models.ChangeEvent is specialised on mongo
	// how to abstract that?
	eventforward.ForwardEvents[models.ChangeEvent](decls)
}
