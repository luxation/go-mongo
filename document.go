package mongo

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Document interface {
	Id() string
	SetId(id uuid.UUID)
	DocumentName() string
	IncrementVersion()
	SetCreatedAt()
	SetUpdatedAt()
}

type BasicDocument struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`
	Version   uint      `json:"version" bson:"version,omitempty"`
}

func (d BasicDocument) Id() string {
	return d.ID
}

func (d *BasicDocument) SetId(id uuid.UUID) {
	d.ID = id.String()
}

func (d *BasicDocument) IncrementVersion() {
	d.Version++
}

func (d *BasicDocument) SetCreatedAt() {
	d.CreatedAt = time.Now()
}

func (d *BasicDocument) SetUpdatedAt() {
	d.UpdatedAt = time.Now()
}

// Deprecated: Use BasicDocument instead.
type MongoDocument struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`
	Version   uint      `json:"version" bson:"version,omitempty"`
}

func (d MongoDocument) Id() string {
	return d.ID
}

func (d *MongoDocument) SetId(id uuid.UUID) {
	d.ID = id.String()
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
