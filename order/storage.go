package order

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	databaseName        = "retail"
	orderCollectionName = "order"
)

type OrderStore interface {
	AddOrders(orders []interface{}) error
	GetOrders(pageSize, pageNumber int) ([]Item, error)
	GetSummaries(pageSize, pageNumber int) ([]Summary, error)
}

type orderDatabase struct {
	client *mongo.Client
}

func NewOrderStore(client *mongo.Client) *orderDatabase {
	return &orderDatabase{
		client: client,
	}
}

func (d *orderDatabase) AddOrders(orders []interface{}) error {
	collection := d.client.Database(databaseName).Collection(orderCollectionName)

	collection.InsertMany(context.Background(), orders)
	return nil
}

func (d *orderDatabase) GetOrders(pageSize, pageNumber int) ([]Item, error) {

	return []Item{}, nil
}

func (d *orderDatabase) GetSummaries(pageSize, pageNumber int) ([]Summary, error) {
	return []Summary{}, nil
}

func (d *orderDatabase) CreateIndexes() error {
	collection := d.client.Database(databaseName).Collection(orderCollectionName)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "customerId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "orderId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "orderId", Value: 1}},
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)

	if err != nil {
		return err
	}

	return nil
}
