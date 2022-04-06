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

type ResultDecoder func(cursor ResultCursor) error

//go:generate mockgen -destination=mocks/client_mock.go -package=mongo . Client
type Client interface {
	Connect() error
	Disconnect() error
	HealthCheck() error
	WithContext(ctx context.Context) Client
	Persist(d Document) error
	GetCollectionByName(name string) (*mongo.Collection, error)
	GetCollection(d Document) (*mongo.Collection, error)
	Aggregate(d Document, pipeline bson.A, decoder ResultDecoder, aggregateOptions ...*options.AggregateOptions) error
	FindAll(d Document, filters bson.M, decoder ResultDecoder, findOptions ...*FindOptions) error
	FindOne(d Document, filters bson.M, findOptions ...*FindOptions) error
	FindOneById(d Document, id string) error
	ReplaceOrPersist(d Document) error
	Replace(d Document) error
	Delete(d Document) error
	DeleteWhere(d Document, key, value string) error
	DeleteMany(d Document, filter bson.M) (int64, error)
	Update(d Document, id string, input interface{}) error
	UpdateWhere(d Document, filter bson.M, input interface{}) error
	UpdateMany(d Document, filter bson.M, input interface{}) error
	GenerateUUID() uuid.UUID
	GetURI() string
}

type mongoClient struct {
	database string
	uri      string
	ctx      *context.Context
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

func (m *mongoClient) WithContext(ctx context.Context) Client {
	m.ctx = &ctx
	return m
}

func (m *mongoClient) GetCollectionByName(name string) (*mongo.Collection, error) {
	if client == nil {
		return nil, errors.New("MongoDB client was not initialized")
	}

	return client.Database(m.database).Collection(name), nil
}

func (m *mongoClient) GetCollection(d Document) (*mongo.Collection, error) {
	if client == nil {
		return nil, errors.New("MongoDB client was not initialized")
	}

	return client.Database(m.database).Collection(d.DocumentName()), nil
}

func (m *mongoClient) Aggregate(d Document, pipeline bson.A, decoder ResultDecoder, opts ...*options.AggregateOptions) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	ag, err := collection.Aggregate(ctx, pipeline, opts...)

	if err != nil {
		if m.ctx != nil {
			m.ctx = nil
		} else {
			cancel()
		}
		return err
	}

	if err = ag.Err(); err != nil {
		if m.ctx != nil {
			m.ctx = nil
		} else {
			cancel()
		}
		return err
	}

	defer ag.Close(ctx)

	for ag.Next(ctx) {
		err = decoder(ResultCursor{Cursor: ag})

		if err != nil {
			if m.ctx != nil {
				m.ctx = nil
			} else {
				cancel()
			}
			return err
		}
	}

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}
	return nil
}

func (m *mongoClient) FindAll(d Document, filters bson.M, decoder ResultDecoder, findOptions ...*FindOptions) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
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
			mongoOptions.Limit = findOption.Pagination.Limit

			if findOption.Pagination.LastID != "" {
				if filters == nil {
					filters = bson.M{
						"_id": bson.M{"$gt": &findOption.Pagination.LastID},
					}
				} else {
					filters["_id"] = bson.M{"$gt": &findOption.Pagination.LastID}
				}
			}
		}
	}

	find, err := collection.Find(ctx, filters, &mongoOptions)

	if err != nil {
		if m.ctx != nil {
			m.ctx = nil
		} else {
			cancel()
		}
		return err
	}

	if err = find.Err(); err != nil {
		if m.ctx != nil {
			m.ctx = nil
		} else {
			cancel()
		}
		return err
	}

	defer find.Close(ctx)

	for find.Next(ctx) {
		err = decoder(ResultCursor{Cursor: find})

		if err != nil {
			if m.ctx != nil {
				m.ctx = nil
			} else {
				cancel()
			}
			return err
		}
	}

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	return nil

}

func (m *mongoClient) FindOne(d Document, filters bson.M, findOptions ...*FindOptions) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	mongoOptions := options.FindOneOptions{}

	if findOptions != nil && findOptions[0] != nil {
		findOption := findOptions[0]

		if findOption.Sort != nil {
			bsonSort := bson.M{}
			for _, sort := range findOption.Sort {
				bsonSort[sort.SortField] = sort.Order
			}
			mongoOptions.Sort = bsonSort
		}
	}

	doc := collection.FindOne(ctx, filters, &mongoOptions)

	if doc.Err() != nil {
		return doc.Err()
	}

	err = doc.Decode(d)

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *mongoClient) FindOneById(d Document, id string) error {
	return m.FindOne(d, bson.M{"_id": id})
}

func (m *mongoClient) Persist(d Document) error {
	if d.GetID() == "" {
		d.SetID(m.GenerateUUID())
	}
	d.SetCreatedAt()
	d.SetUpdatedAt()

	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	_, err = collection.InsertOne(ctx, d)

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	return err
}

func (m *mongoClient) ReplaceOrPersist(d Document) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

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

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	return nil
}

func (m *mongoClient) Replace(d Document) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	d.SetUpdatedAt()

	filter := bson.M{"_id": d.GetID()}

	err = collection.FindOneAndReplace(ctx, filter, d).Err()

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	return err
}

func (m *mongoClient) Delete(d Document) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	filter := bson.M{"_id": d.GetID()}

	dr, err := collection.DeleteOne(ctx, filter)

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

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

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	filter := bson.M{key: value}

	dr, err := collection.DeleteOne(ctx, filter)

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	if err != nil {
		return err
	}

	if dr.DeletedCount != 1 {
		return errors.New(fmt.Sprintf("Deleted %q elements", dr.DeletedCount))
	}

	return nil
}

func (m *mongoClient) DeleteMany(d Document, filter bson.M) (int64, error) {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return 0, err
	}

	if collection == nil {
		return 0, errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	dr, err := collection.DeleteMany(ctx, filter)

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	if err != nil {
		return 0, err
	}

	return dr.DeletedCount, nil
}

func (m *mongoClient) Update(d Document, id string, input interface{}) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

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

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	return err
}

func (m *mongoClient) UpdateMany(d Document, filter bson.M, input interface{}) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

	collection, err := m.GetCollection(d)

	if err != nil {
		return err
	}

	if collection == nil {
		return errors.New(fmt.Sprintf("No collection found for document named %s", d.DocumentName()))
	}

	updates := FlattenedMapFromInterface(input)
	updates["updatedAt"] = time.Now()

	_, err = collection.UpdateMany(ctx, filter, bson.D{
		{Key: "$set", Value: updates},
	})

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

	return err
}

func (m *mongoClient) UpdateWhere(d Document, filter bson.M, input interface{}) error {
	ctx, cancel := m.getContext()

	if m.ctx != nil {
		ctx = *m.ctx
		cancel()
	}

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

	if m.ctx != nil {
		m.ctx = nil
	} else {
		cancel()
	}

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
