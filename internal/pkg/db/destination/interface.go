package destination

type DB[T any] interface {
	SendOperations(done chan struct{}, opChan <-chan *T, errChan chan<- error, database, table string)
}
