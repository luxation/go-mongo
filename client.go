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

type ResultCursor struct {
	*mongo.Cursor
}

type Client interface {
	Connect() error
	Disconnect() error
	HealthCheck() error
	Persist(d Document) error
	PersistWithContext(d Document, ctx context.Context) error
	GetCollection(d Document) (*mongo.Collection, error)
	FindAll(d Document, filters bson.M, findOptions ...*FindOptions) (*ResultCursor, error)
	FindAllWithContext(d Document, filters bson.M, ctx context.Context, findOptions ...*FindOptions) (*ResultCursor, error)
	FindOne(d Document, filters bson.M) error
	FindOneWithContext(d Document, filters bson.M, ctx context.Context) error
	FindOneById(d Document, id string) error
	FindOneByIdWithContext(d Document, id string, ctx context.Context) error
	ReplaceOrPersist(d Document) error
	ReplaceOrPersistWithContext(d Document, ctx context.Context) error
	Replace(d Document) error
	ReplaceWithContext(d Document, ctx context.Context) error
	Delete(d Document) error
	DeleteWhere(d Document, key, value string) error
	DeleteWithContext(d Document, ctx context.Context) error
	Update(d Document, id string, input interface{}) error
	UpdateWhere(d Document, filter bson.M, input interface{}) error
	UpdateWithContext(d Document, id string, input interface{}, ctx context.Context) error
	GenerateUUID() uuid.UUID
	GetURI() string
}

type mongoClient struct {
	database string
	uri      string
}

func (m *mongoClient) FindAll(d Document, filters bson.M, findOptions ...*FindOptions) (*ResultCursor, error) {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return nil, err
	}

	mongoOptions := options.FindOptions{}

	if findOptions != nil && findOptions[0] != nil {
		findOption := findOptions[0]

		if findOption.Sort != nil {
			bsonSort := bson.M{}
			for _, sort := range findOption.Sort {
				bsonSort[sort.SortField] = sort.Order
			}
			mongoOptions.Sort = bsonSort
		}

		if findOption.Pagination != nil {
			mongoOptions.Limit = &findOption.Pagination.Limit
			skip := findOption.Pagination.Page * findOption.Pagination.Limit
			mongoOptions.Skip = &skip
		}
	}

	find, err := collection.Find(ctx, filters, &mongoOptions)
	if err != nil {
		return nil, err
	}

	return &ResultCursor{Cursor: find}, nil

}

func (m *mongoClient) FindAllWithContext(d Document, filters bson.M, ctx context.Context, findOptions ...*FindOptions) (*ResultCursor, error) {
	collection, err := m.GetCollection(d)

	if err != nil {
		return nil, err
	}

	mongoOptions := options.FindOptions{}

	if findOptions != nil && findOptions[0] != nil {
		findOption := findOptions[0]

		if findOption.Sort != nil {
			bsonSort := bson.M{}
			for _, sort := range findOption.Sort {
				bsonSort[sort.SortField] = sort.Order
			}
			mongoOptions.Sort = bsonSort
		}

		if findOption.Pagination != nil {
			mongoOptions.Limit = &findOption.Pagination.Limit
			skip := findOption.Pagination.Page * findOption.Pagination.Limit
			mongoOptions.Skip = &skip
		}
	}

	find, err := collection.Find(ctx, filters, &mongoOptions)

	if err != nil {
		return nil, err
	}

	return &ResultCursor{Cursor: find}, nil
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
	if d.GetID() == "" {
		d.SetID(m.GenerateUUID())
	}
	d.SetCreatedAt()
	d.SetUpdatedAt()

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

func (m *mongoClient) PersistWithContext(d Document, ctx context.Context) error {
	if d.GetID() == "" {
		d.SetID(m.GenerateUUID())
	}
	d.SetCreatedAt()
	d.SetUpdatedAt()

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

func (m *mongoClient) FindOneWithContext(d Document, filters bson.M, ctx context.Context) error {
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

func (m *mongoClient) FindOneByIdWithContext(d Document, id string, ctx context.Context) error {
	return m.FindOneWithContext(d, bson.M{"_id": id}, ctx)
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

	d.SetCreatedAt()
	d.SetUpdatedAt()

	filter := bson.M{"_id": d.GetID()}

	err = collection.FindOneAndReplace(ctx, filter, d).Err()

	if err != nil {
		return m.Persist(d)
	}

	return nil
}

func (m *mongoClient) ReplaceOrPersistWithContext(d Document, ctx context.Context) error {
	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	d.SetCreatedAt()
	d.SetUpdatedAt()

	filter := bson.M{"_id": d.GetID()}

	err = collection.FindOneAndReplace(ctx, filter, d).Err()

	if err != nil {
		return m.PersistWithContext(d, ctx)
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

	d.SetUpdatedAt()

	filter := bson.M{"_id": d.GetID()}

	return collection.FindOneAndReplace(ctx, filter, d).Err()
}

func (m *mongoClient) ReplaceWithContext(d Document, ctx context.Context) error {
	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	d.SetUpdatedAt()

	filter := bson.M{"_id": d.GetID()}

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

	filter := bson.M{"_id": d.GetID()}

	dr, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if dr.DeletedCount != 1 {
		return errors.New(fmt.Sprintf("Deleted %q elements", dr.DeletedCount))
	}

	return nil
}

func (m *mongoClient) DeleteWhere(d Document, key, value string) error {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	filter := bson.M{key: value}

	dr, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if dr.DeletedCount != 1 {
		return errors.New(fmt.Sprintf("Deleted %q elements", dr.DeletedCount))
	}

	return nil
}

func (m *mongoClient) DeleteWithContext(d Document, ctx context.Context) error {
	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	filter := bson.M{"_id": d.GetID()}

	dr, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if dr.DeletedCount != 1 {
		return errors.New(fmt.Sprintf("Deleted %q elements", dr.DeletedCount))
	}

	return nil
}

func (m *mongoClient) Update(d Document, id string, input interface{}) error {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	filter := bson.M{"_id": id}

	updates := FlattenedMapFromInterface(input)
	updates["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updates},
	})

	return err
}

func (m *mongoClient) UpdateWhere(d Document, filter bson.M, input interface{}) error {
	ctx, cancel := m.getContext()
	defer cancel()

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	updates := FlattenedMapFromInterface(input)
	updates["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updates},
	})

	return err
}

func (m *mongoClient) UpdateWithContext(d Document, id string, input interface{}, ctx context.Context) error {
	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	filter := bson.M{"_id": id}

	updates := FlattenedMapFromInterface(input)
	updates["updatedAt"] = time.Now()

	_, err = collection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updates},
	})

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
