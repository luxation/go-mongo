package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

var (
	client *mongo.Client
)

type Client interface {
	getContext() (context.Context, context.CancelFunc)
	Connect() error
	Disconnect() error
	HealthCheck() error
	Persist(d Document) error
	GetCollection(d Document) (*mongo.Collection, error)
	FindOne(d Document, filters bson.M) error
	FindOneById(d Document, id string) error
	ReplaceOrPersist(d Document) error
	Replace(d Document) error
	Delete(d Document) error
	Update(d Document, id string) error
	GenerateUUID() uuid.UUID
	GetURI() string
}

type mongoClient struct {
	database string
	uri      string
}

func (m *mongoClient) getContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

func (m *mongoClient) Connect() error {
	ctx, cancel := m.getContext()
	defer cancel()

	c, err := mongo.Connect(ctx, options.Client().ApplyURI(m.uri))

	if err != nil {
		return err
	}

	client = c

	return nil
}

func (m *mongoClient) Disconnect() error {
	if client == nil {
		return errors.New("MongoDB client was not initialized")
	}

	ctx, cancel := m.getContext()
	defer cancel()
	return client.Disconnect(ctx)
}

func (m *mongoClient) HealthCheck() error {
	if client == nil {
		return errors.New("MongoDB client was not initialized")
	}

	ctx, cancel := m.getContext()
	defer cancel()

	return client.Ping(ctx, readpref.Primary())
}

func (m *mongoClient) Persist(d Document) error {
	if d.Id() == "" {
		d.SetId(m.GenerateUUID())
	}
	d.SetCreatedAt()
	d.SetUpdatedAt()
	d.IncrementVersion()

	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	_, err = collection.InsertOne(ctx, d)

	return err
}

func (m *mongoClient) GetCollection(d Document) (*mongo.Collection, error) {
	if client == nil {
		return nil, errors.New("MongoDB client was not initialized")
	}

	return client.Database(m.database).Collection(d.DocumentName()), nil
}

func (m *mongoClient) FindOne(d Document, filters bson.M) error {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	doc := collection.FindOne(ctx, filters)

	if doc.Err() != nil {
		return doc.Err()
	}

	err = doc.Decode(d)

	if err != nil {
		return err
	}

	return nil
}

func (m *mongoClient) FindOneById(d Document, id string) error {
	return m.FindOne(d, bson.M{"_id": id})
}

func (m *mongoClient) ReplaceOrPersist(d Document) error {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	d.IncrementVersion()
	d.SetCreatedAt()
	d.SetUpdatedAt()

	filter := bson.M{"_id": d.Id()}

	err = collection.FindOneAndReplace(ctx, filter, d).Err()

	if err != nil {
		return m.Persist(d)
	}

	return nil
}

func (m *mongoClient) Replace(d Document) error {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	d.IncrementVersion()
	d.SetUpdatedAt()

	filter := bson.M{"_id": d.Id()}

	return collection.FindOneAndReplace(ctx, filter, d).Err()
}

func (m *mongoClient) Delete(d Document) error {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	filter := bson.M{"_id": d.Id()}

	dr, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if dr.DeletedCount != 1 {
		return errors.New(fmt.Sprintf("Deleted %q elements", dr.DeletedCount))
	}

	return nil
}

func (m *mongoClient) Update(d Document, id string) error {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	d.IncrementVersion()
	d.SetUpdatedAt()

	filter := bson.M{"_id": id}

	_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": d})

	return err
}

func (m *mongoClient) GenerateUUID() uuid.UUID {
	return uuid.New()
}

func (m *mongoClient) GetURI() string {
	return m.uri
}

func NewClient(config ClientConfig) (Client, error) {
	newClient := &mongoClient{
		database: config.Database,
	}

	uri, err := config.generateURI()

	if err != nil {
		return nil, err
	}

	newClient.uri = uri

	err = newClient.Connect()

	if err != nil {
		return nil, err
	}

	return newClient, nil
}

// Deprecated: Use NewClient instead.
func NewMongoClient(config ClientConfig) (Client, error) {
	newClient := &mongoClient{
		database: config.Database,
	}

	uri, err := config.generateURI()

	if err != nil {
		return nil, err
	}

	newClient.uri = uri

	err = newClient.Connect()

	if err != nil {
		return nil, err
	}

	return newClient, nil
}
