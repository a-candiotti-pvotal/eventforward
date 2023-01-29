package source

import (
	"eventforward/internal/pkg/models"
)

type DB interface {
	WatchOperations(done chan struct{}, opChan chan<- *models.ChangeEvent, errChan chan<- error, from string)
	ReadOperations(done chan struct{}, opChan chan<- *models.ChangeEvent, errChan chan<- error, namespace string)
}
