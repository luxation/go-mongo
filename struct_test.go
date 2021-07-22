package mongo

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

type FooDocument struct {
	BasicDocument
	Action string
}

func (f FooDocument) DocumentName() string {
	return "foo"
}

func (f *FooDocument) FromBSON(sr *mongo.SingleResult) error {
	err := sr.Decode(&f)

	if err != nil {
		return err
	}

	return nil
}

func TestCorrectInheritedType(t *testing.T) {
	uuidField := uuid.NewString()
	fooDoc := &FooDocument{
		BasicDocument: BasicDocument{
			ID:        uuidField,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Action: "Bar",
	}

	assert.Equal(t, "foo", fooDoc.DocumentName())
	assert.Equal(t, uuidField, fooDoc.GetID())

	fooDoc.IncrementVersion()

	assert.Equal(t, 1, fooDoc.Version)
}
