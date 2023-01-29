package mongo

import (
	"os"
	"testing"
	"time"

	"eventforward/internal/pkg/models"
)

var m *MongoDB

func injectAndReadBack(done chan struct{}, opChan chan *models.ChangeEvent, total int) {
	go func() {
		time.Sleep(1 * time.Second)
		m.injectDataIn("profiling", "bet", total)
	}()

	count := 0
	go func() {
		for {
			select {
			case <-opChan:
				count += 1
				if count == total {
					close(done)
				}
			}
		}
	}()
}

func init() {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	var err error
	m, err = Setup()
	if err != nil {
		panic(err)
	}
}

func BenchmarkRecvOperations(b *testing.B) {
	for _, total := range []int{1, 1000, 100000} {
		var done = make(chan struct{})
		var opChan = make(chan *models.ChangeEvent, 5000)
		var errChan = make(chan error) // Error Channel

		injectAndReadBack(done, opChan, total)
		m.RecvOperations(done, opChan, errChan)
	}
}
