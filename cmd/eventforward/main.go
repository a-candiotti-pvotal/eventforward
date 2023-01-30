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
	{From: "profiling.bet", To: "bet", Watch: true},
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// FIXME : models.ChangeEvent is specialised on mongo
	// how to abstract that?
	eventforward.ForwardEvents[models.ChangeEvent](decls)
}
