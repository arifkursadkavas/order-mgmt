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
	AddOrders(ctx context.Context, orders []interface{}) error
	GetOrders(ctx context.Context, pageSize, pageNumber int) ([]Item, error)
	GetSummaries(ctx context.Context, pageSize, pageNumber int) ([]Summary, error)
}

type orderDatabase struct {
	client *mongo.Client
}

func NewOrderStore(client *mongo.Client) *orderDatabase {
	return &orderDatabase{
		client: client,
	}
}

func (d *orderDatabase) AddOrders(ctx context.Context, orders []interface{}) error {
	collection := d.client.Database(databaseName).Collection(orderCollectionName)

	collection.InsertMany(context.Background(), orders)
	return nil
}

func (d *orderDatabase) GetOrders(ctx context.Context, pageSize, pageNumber int) ([]Item, error) {

	collection := d.client.Database(databaseName).Collection(orderCollectionName)

	cursor, err := collection.Find(ctx, bson.D{}, nil)

	if err != nil {
		return []Item{}, nil
	}

	var items []Item
	if err = cursor.All(ctx, &items); err != nil {
		panic(err)
	}

	return items, err
}

func (d *orderDatabase) GetSummaries(ctx context.Context, pageSize, pageNumber int) ([]Summary, error) {
	collection := d.client.Database(databaseName).Collection(orderCollectionName)

	groupStage := bson.D{
		{"$group",
			bson.D{
				{"_id", "$customerId"},
				{"customerId", bson.D{{"$first", "$customerId"}}},
				{"nbrOfPurchasedItems", bson.D{{"$sum", 1}}},
				{"totalAmountEur", bson.D{{"$sum", "$costEur"}}},
			},
		},
	}

	pipeline := mongo.Pipeline{groupStage}

	cursor, err := collection.Aggregate(ctx, pipeline, nil)

	if err != nil {
		return []Summary{}, nil
	}

	var items []Summary
	if err = cursor.All(ctx, &items); err != nil {
		panic(err)
	}

	return items, err
}

func (d *orderDatabase) CreateIndexes() error {
	collection := d.client.Database(databaseName).Collection(orderCollectionName)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "customerId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "itemId", Value: 1}},
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
