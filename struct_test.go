package mongo

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"time"
)

type FooDocument struct {
	MongoDocument
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
	fooDoc := &FooDocument{
		MongoDocument: MongoDocument{
			ID:        "123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   0,
		},
		Action: "Bar",
	}

	assert.Equal(t, "foo", fooDoc.DocumentName())
	assert.True(t, fooDoc.IsUniqueID())
	assert.Equal(t, "123", fooDoc.Id())

	fooDoc.IncrementVersion()

	assert.Equal(t, uint(1), fooDoc.Version)
}
