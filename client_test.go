package mongo

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

type Foo struct {
	MongoDocument `bson:",inline"`
	Action        string
}

func (f Foo) Id() string { return f.ID }

func (f Foo) DocumentName() string { return "foo" }

func (f Foo) IsUniqueID() bool {
	return true
}

func (f Foo) FromBSON(sr *mongo.SingleResult) (Document, error) {
	err := sr.Decode(&f)

	if err != nil {
		return nil, err
	}

	return &f, nil
}

func dummyConnect() (Client, error) {
	clientConfig := ClientConfig{
		Host:     "localhost",
		Port:     27017,
		Database: "test_db",
	}

	return NewMongoClient(clientConfig)
}

// WARNING: Those tests should be run in order or at least TestConnect should
// be launched before any other test to initialize the client

var (
	testClient   Client
	connectError error
)

func TestConnect(t *testing.T) {
	testClient, connectError = dummyConnect()

	assert.Nil(t, connectError)
	assert.NotNil(t, testClient)
}

func TestInsertFoo(t *testing.T) {
	assert.NotNil(t, testClient)

	foo := Foo{
		MongoDocument: MongoDocument{
			ID: "foo-bar-1",
		},
		Action: "Bar",
	}

	err := testClient.Persist(&foo)

	assert.Nil(t, err)

	err = testClient.Persist(&foo)

	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("Document %s with ID %s already exists", foo.DocumentName(), foo.Id()), err.Error())
}

func TestFindOneByID(t *testing.T) {
	assert.NotNil(t, testClient)

	result, err := testClient.FindOneById(&Foo{}, "foo-bar-1")

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "foo-bar-1", result.Id())
}

func TestReplaceOrPersistReplace(t *testing.T) {
	assert.NotNil(t, testClient)

	foo := Foo{
		Action: "Bar Replaced",
		MongoDocument: MongoDocument{
			ID: "foo-bar-1",
		},
	}

	err := testClient.ReplaceOrPersist(&foo)

	assert.Nil(t, err)
}

func TestReplaceOrPersistPersist(t *testing.T) {
	assert.NotNil(t, testClient)

	foo := Foo{
		Action: "Bar Persisted",
		MongoDocument: MongoDocument{
			ID: "foo-bar-2",
		},
	}

	err := testClient.ReplaceOrPersist(&foo)

	assert.Nil(t, err)
}

func TestFindOneByIDNewlyPersisted(t *testing.T) {
	assert.NotNil(t, testClient)

	result, err := testClient.FindOneById(&Foo{}, "foo-bar-2")

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "foo-bar-2", result.Id())
}

func TestDelete(t *testing.T) {
	assert.NotNil(t, testClient)

	foo := Foo{
		Action: "Bar Persisted",
		MongoDocument: MongoDocument{
			ID: "foo-bar-2",
		},
	}

	err := testClient.Delete(&foo)

	assert.Nil(t, err)
}

func TestUpdate(t *testing.T) {
	assert.NotNil(t, testClient)

	result, err := testClient.FindOneById(&Foo{}, "foo-bar-1")

	assert.Nil(t, err)

	result.IncrementVersion()
	result.SetUpdatedAt()

	err = testClient.Update(result, "foo-bar-1", bson.M{"action": "updated bar with update method"})

	assert.Nil(t, err)
}
