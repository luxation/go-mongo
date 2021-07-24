package mongo

import (
	"github.com/google/uuid"
	"time"
)

type Document interface {
	GetID() string
	SetID(id uuid.UUID)
	DocumentName() string
	SetCreatedAt()
	SetUpdatedAt()
}

type BasicDocument struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}

func (d BasicDocument) GetID() string {
	return d.ID
}

func (d *BasicDocument) SetID(id uuid.UUID) {
	d.ID = id.String()
}

func (d *BasicDocument) SetCreatedAt() {
	d.CreatedAt = time.Now()
}

func (d *BasicDocument) SetUpdatedAt() {
	d.UpdatedAt = time.Now()
}
