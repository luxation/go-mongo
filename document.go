package mongo

import (
	"github.com/google/uuid"
	"time"
)

type Document interface {
	GetID() string
	SetID(id uuid.UUID)
	DocumentName() string
	IncrementVersion()
	SetCreatedAt()
	SetUpdatedAt()
}

type BasicDocument struct {
	ID        string    `json:"id" bson:"_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

func (d BasicDocument) GetID() string {
	return d.ID
}

func (d *BasicDocument) SetID(id uuid.UUID) {
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
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   uint      `json:"version"`
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
