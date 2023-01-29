package destination

import (
	"eventforward/internal/pkg/models"
)

type DB interface {
	SendOperations(done chan struct{}, opChan <-chan *models.ChangeEvent, errChan chan<- error, to string)
}
