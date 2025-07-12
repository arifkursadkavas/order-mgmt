package order

import (
	"context"

	"company.com/retail/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderStore interface {
	AddOrders(orders []interface{}) error
	GetOrders(pageSize, pageNumber int) ([]Item, error)
	GetSummaries(pageSize, pageNumber int) ([]Summary, error)
}

type orderDatabase struct {
	client *mongo.Client
}

func NewOrderStore(client *mongo.Client) OrderStore {
	return &orderDatabase{
		client: client,
	}
}

func (d *orderDatabase) AddOrders(orders []interface{}) error {
	collection := config.Config.DBClient.Database("retail").Collection("order")

	collection.InsertMany(context.Background(), orders)
	return nil
}

func (d *orderDatabase) GetOrders(pageSize, pageNumber int) ([]Item, error) {

	return []Item{}, nil
}

func (d *orderDatabase) GetSummaries(pageSize, pageNumber int) ([]Summary, error) {
	return []Summary{}, nil
}
