package order

import "go.mongodb.org/mongo-driver/mongo"

type OrderStore interface {
	AddOrders(orders []OrderDb) error
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

func (d *orderDatabase) AddOrders(orders []OrderDb) error {

	return nil
}

func (d *orderDatabase) GetOrders(pageSize, pageNumber int) ([]Item, error) {

	return []Item{}, nil
}

func (d *orderDatabase) GetSummaries(pageSize, pageNumber int) ([]Summary, error) {
	return []Summary{}, nil
}
