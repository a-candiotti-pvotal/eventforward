package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"eventforward/internal/pkg/db/destination"
	"eventforward/internal/pkg/db/source"
	"eventforward/internal/pkg/models"
)

type ForwardDecl struct {
	From string
	To string
	Watch bool
}

func forward(done chan struct{}, wg *sync.WaitGroup, src source.DB, dst destination.DB, decl *ForwardDecl) {
	opChan := make(chan *models.ChangeEvent, 5000)
	errChan := make(chan error) // Error Channel

	// FIXME : cleaner wg?
	//         what to do with errors?
	go func () {
		dst.SendOperations(done, opChan, errChan, decl.To)
		wg.Done()
	}()

	if decl.Watch {
		src.WatchOperations(done, opChan, errChan, decl.From)
	} else {
		// TODO : rename?
//		src.ReadOperations(done, opChan, errChan, decl.From)
	}
	wg.Done()
}

// FIXME : read from yaml/json
var decls = []ForwardDecl{
	{From: "profiling.bet", To: "bet"},
}

func forwardEvents() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	go func() {
		sig := <-sigs
		println()
		fmt.Println(sig)
		close(done)
	}()

	srcdb := source.DBFromEnv()
	dstdb := destination.DBFromEnv()

	wg := &sync.WaitGroup{}
	wg.Add(len(decls) * 2)

	// TODO : add flag to replay instead of forward

	for _, decl := range decls {
		go forward(done, wg, srcdb, dstdb, &decl)
	}

	wg.Wait()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	forwardEvents()
}
