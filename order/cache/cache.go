package cache

import (
	"errors"
	"time"

	"company.com/order-service/order/model"
)

var (
	ordersCacheKey        = "orders"
	ordersSummaryCacheKey = "summaries"
)

//go:generate mockery --name OrderCache --output ../../internal/mocks --with-expecter
type OrderCache interface {
	AddOrders(orders []model.OrderCacheModel, orderSummaries map[string]model.OrderSummaryCacheModel) error
	GetOrders() ([]model.Item, error)
	GetSummaries() ([]model.Summary, error)
}

//go:generate mockery --name Cache --output ../../internal/mocks --with-expecter
type Cache interface {
	Get(k string) (interface{}, bool)
	Add(k string, x interface{}, d time.Duration) error
	Set(k string, x interface{}, d time.Duration)
}

type orderCache struct {
	cache             Cache
	defaultExpiration time.Duration
}

func NewOrderCache(cache Cache, defaultExpiration time.Duration) *orderCache {

	return &orderCache{
		cache:             cache,
		defaultExpiration: defaultExpiration,
	}
}

// This method adds the orders in the format expected by the two listing methods.
// One cache for the listing individual items, other for the summaries
func (c *orderCache) AddOrders(orders []model.OrderCacheModel, incomingSummaries map[string]model.OrderSummaryCacheModel) error {

	//Individual items with order data

	//check if orders cache already exists
	o, e := c.cache.Get(ordersCacheKey)

	//if exists update the existing
	if e {
		existingOrders := o.([]model.OrderCacheModel)
		existingOrders = append(existingOrders, orders...)
		c.cache.Set(ordersCacheKey, existingOrders, c.defaultExpiration)
	} else {
		//else add for the first time
		err := c.cache.Add(ordersCacheKey, orders, c.defaultExpiration)

		if err != nil {
			return err
		}
	}

	//Summaries

	//Check if order summaries cache already exists
	orderSums, exists := c.cache.Get(ordersSummaryCacheKey)

	if exists {
		//if exists iterate through each incoming order to check if there exists a customerId key in the existing cache
		existingOrderSummaries := orderSums.(map[string]model.OrderSummaryCacheModel)

		for customerId, os := range incomingSummaries {
			if o, ok := (existingOrderSummaries)[customerId]; ok {
				var update = model.OrderSummaryCacheModel{
					CustomerId:          customerId,
					NbrOfPurchasedItems: o.NbrOfPurchasedItems + os.NbrOfPurchasedItems,
					TotalAmountEur:      o.TotalAmountEur + os.TotalAmountEur,
				}
				(existingOrderSummaries)[customerId] = update
			}
		}
		c.cache.Set(ordersSummaryCacheKey, existingOrderSummaries, c.defaultExpiration)
	} else { //If this is the first time summary cache is being written, write incoming directly.
		c.cache.Add(ordersSummaryCacheKey, incomingSummaries, c.defaultExpiration)
	}

	return nil
}

// Reviewcomments --- no filtering/querying by customerId
func (c *orderCache) GetOrders() ([]model.Item, error) {

	o, found := c.cache.Get(ordersCacheKey)

	if !found {
		return []model.Item{}, errors.New("no order found")
	}

	orders := o.([]model.OrderCacheModel)

	var items []model.Item
	for _, o := range orders {
		items = append(items, model.Item{
			CustomerId: o.CustomerId,
			ItemId:     o.ItemId,
			CostEur:    o.CostEur,
		})
	}

	return items, nil
}

func (c *orderCache) GetSummaries() ([]model.Summary, error) {

	o, found := c.cache.Get(ordersSummaryCacheKey)

	if !found {
		return []model.Summary{}, errors.New("no summary found")
	}

	orderSummaries := o.(map[string]model.OrderSummaryCacheModel)

	var items []model.Summary
	for _, o := range orderSummaries {
		items = append(items, model.Summary{
			CustomerId:          o.CustomerId,
			NbrOfPurchasedItems: o.NbrOfPurchasedItems,
			TotalAmountEur:      o.TotalAmountEur,
		})
	}

	return items, nil
}
