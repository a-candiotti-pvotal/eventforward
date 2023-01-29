package models

import (
	"time"
)

type documentKey struct {
	ID string `bson:"_id"`
}

type changeID struct {
	Data string `bson:"_data"`
}

type namespace struct {
	Db   string `bson:"db"`
	Coll string `bson:"coll"`
}

// TODO : remove mongo specific fields as much as possible
type ChangeEvent struct {
	ID            changeID            		`bson:"_id"`
	OperationType string              		`bson:"operationType"`
	ClusterTime   time.Time 							`bson:"clusterTime"`
	FullDocument  map[string]interface{} 	`bson:"fullDocument"`
	DocumentKey   documentKey         		`bson:"documentKey"`
	Ns            namespace           		`bson:"ns"`
}
