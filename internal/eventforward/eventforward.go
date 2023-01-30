package eventforward

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"eventforward/pkg/models"
	"eventforward/internal/pkg/db/destination"
	"eventforward/internal/pkg/db/source"
)

func forward[T any](done chan struct{}, wg *sync.WaitGroup, decl *models.ForwardDecl) {
	src := source.DBFromName[T](decl.Name, decl.From.Type)
	dst := destination.DBFromName[T](decl.Name, decl.To.Type)

	opChan := make(chan *T, 5000)
	errChan := make(chan error) // Error Channel

	// FIXME : cleaner wg?
	//         what to do with errors?
	go func () {
		dst.SendOperations(done, opChan, errChan, decl.To.Database, decl.To.Table)
		wg.Done()
	}()

	if decl.Watch {
		src.WatchOperations(done, opChan, errChan, decl.From.Database, decl.To.Table)
	} else {
		// TODO : rename?
		src.ReadOperations(done, opChan, errChan, decl.From.Database, decl.From.Table)
	}
	wg.Done()
}

func ForwardEvents[T any](decls []models.ForwardDecl) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		sig := <-sigs
		println()
		fmt.Println(sig)
		close(done)
	}()

	wg := &sync.WaitGroup{}
	wg.Add(len(decls) * 2)

	for _, decl := range decls {
		if decl.From.Database == "any" {
			go forward[map[string]interface{}](done, wg, &decl)
		} else {
			go forward[T](done, wg, &decl)
		}
	}

	wg.Wait()
}
