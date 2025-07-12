package order

import "go.mongodb.org/mongo-driver/mongo"

type OrderDbI interface {
	AddOrder() error
	GetOrders() error
	GetSummaries() error
}

func NewOrderDatabase(client *mongo.Client) *OrderDbI {

}
