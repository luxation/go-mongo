package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Document interface {
	Id() string
	SetId(id string)
	DocumentName() string
	FromBSON(sr *mongo.SingleResult) (Document, error)
	IsUniqueID() bool
	IncrementVersion()
	SetCreatedAt()
	SetUpdatedAt()
}

type MongoDocument struct {
	ID        string    `json:"id" bson:"id"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
	Version   uint      `json:"version" bson:"version"`
}

func (d MongoDocument) Id() string {
	return d.ID
}

func (d *MongoDocument) SetId(id string) {
	d.ID = id
}

func (d MongoDocument) IsUniqueID() bool {
	return true
}

func (d *MongoDocument) IncrementVersion() {
	d.Version++
}

func (d *MongoDocument) SetCreatedAt() {
	d.CreatedAt = time.Now()
}

func (d *MongoDocument) SetUpdatedAt() {
	d.UpdatedAt = time.Now()
}

func ToBSON(d Document) ([]byte, error) {
	bsonRes, err := bson.Marshal(d)

	if err != nil {
		return nil, err
	}

	return bsonRes, nil
}

func ToDoc(v interface{}) (doc *bson.M, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return nil, err
	}

	err = bson.Unmarshal(data, &doc)
	return doc, nil
}
