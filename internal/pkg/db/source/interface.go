package source

type DB[T any] interface {
	WatchOperations(done chan struct{}, opChan chan<- *T, errChan chan<- error, database, table string)
	ReadOperations(done chan struct{}, opChan chan<- *T, errChan chan<- error, database, table string)
}
