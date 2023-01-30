package mongo

import (
	"context"
	"log"
	"time"
)

func createNDocs(ndocs int) []interface{} {
	results := []interface{}{}
	for ndocs > 0 {
		results = append(results, map[string]string{
			"key": "value",
		})
		ndocs -= 1
	}

	return results
}

func (m *MongoDB[T]) injectDataIn(targetDatabaseName, targetCollectionName string, nRecords int) {
	targetCollection := m.client.Database(targetDatabaseName).Collection(targetCollectionName)

	docs := createNDocs(nRecords)

	//			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ctx := context.Background()
	_, err := targetCollection.InsertMany(ctx, docs)
	if err != nil {
		log.Println(err)
		//				cancel()
		return
	}
//	log.Println("DEBUG : All records were injected")
}

// TODO : random ticker and random number of injected records
// https://w11i.me/random-ticker-in-go
func (m *MongoDB[T]) injectDataInEvery(done chan struct{}, interval time.Duration, targetDatabaseName, targetCollectionName string, nRecords int) {
//	log.Printf("DEBUG : injecting %d records in %s every %s\n", nRecords, targetCollectionName, interval)

	ticker := time.NewTicker(interval)

	index := 0
	total := nRecords

	for index < total {
		select {
		case <-done:
			return

		case <-ticker.C:
			m.injectDataIn(targetDatabaseName, targetCollectionName, nRecords)
			index += nRecords
//			log.Printf("DEBUG : injected %d/%d records\n", index, total)
		}
	}

//	log.Println("DEBUG : All records were injected")
}
