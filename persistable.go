package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Document interface {
	Id() string
	EntityName() string
	FromBson(sr *mongo.SingleResult) (Document, error)
	IsUniqueID() bool
}

// ToBSON Convert document object to BSON
func ToBSON(d Document) ([]byte, error) {
	bsonRes, err := bson.Marshal(d)

	if err != nil {
		return nil, err
	}

	return bsonRes, nil
}
