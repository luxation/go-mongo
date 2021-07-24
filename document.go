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
	ID        string    `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at,omitempty"`
	Version   int       `json:"version" bson:"version,omitempty"`
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
