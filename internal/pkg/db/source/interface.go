package source

type DB[T any] interface {
	WatchOperations(done chan struct{}, opChan chan<- *T, errChan chan<- error, from string)
	ReadOperations(done chan struct{}, opChan chan<- *T, errChan chan<- error, namespace string)
}
