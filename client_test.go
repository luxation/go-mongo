package mongo

import (
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type Foo struct {
	Action  string
}

func (f Foo) Id() string { return "foo-id" }

func (f Foo) EntityName() string { return "foo" }

func (f Foo) FromBson(sr *mongo.SingleResult) (Document, error) {
	err := sr.Decode(&f)

	if err != nil {
		return nil, err
	}

	return f, nil
}

func dummyConnect() (Client, error) {
	clientConfig := ClientConfig{
		Host:     "localhost",
		Port:     27017,
		Database: "test_db",
	}

	return NewMongoClient(clientConfig)
}

func TestConnect(t *testing.T) {
	testClient, err := dummyConnect()

	assert.Nil(t, err)
	assert.NotNil(t, testClient)
}

func TestInsertFoo(t *testing.T) {
	testClient, err := dummyConnect()

	assert.Nil(t, err)
	assert.NotNil(t, testClient)

	foo := Foo{
		Action:  "Bar",
	}

	err = testClient.Persist(foo)

	assert.Nil(t, err)
}
