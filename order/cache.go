package order

import (
	"errors"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	ordersCacheKey        = "orders"
	ordersSummaryCacheKey = "summaries"
)

//go:generate mockery --name OrderCache --output ../internal/mocks --with-expecter
type OrderCache interface {
	AddOrders(orders []OrderCacheModel, orderSummaries map[string]OrderSummaryCacheModel) error
	GetOrders() ([]Item, error)
	GetSummaries() ([]Summary, error)
}

type orderCache struct {
	cache             *cache.Cache
	defaultExpiration time.Duration
}

func NewOrderCache(cache *cache.Cache, defaultExpiration time.Duration) *orderCache {

	return &orderCache{
		cache:             cache,
		defaultExpiration: defaultExpiration,
	}
}

// This method adds the orders in the format expected by the two listing methods.
// One cache for the listing individual items, other for the summaries
func (c *orderCache) AddOrders(orders []OrderCacheModel, incomingSummaries map[string]OrderSummaryCacheModel) error {

	//Individual items with order data

	//check if orders cache already exists
	o, e := c.cache.Get(ordersCacheKey)

	//if exists update the existing
	if e {
		existingOrders := o.([]OrderCacheModel)
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
		existingOrderSummaries := orderSums.(map[string]OrderSummaryCacheModel)

		for customerId, os := range incomingSummaries {
			if o, ok := (existingOrderSummaries)[customerId]; ok {
				var update = OrderSummaryCacheModel{
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

func (c *orderCache) GetOrders() ([]Item, error) {

	o, found := c.cache.Get(ordersCacheKey)

	if !found {
		return []Item{}, errors.New("no order found")
	}

	orders := o.([]OrderCacheModel)

	var items []Item
	for _, o := range orders {
		items = append(items, Item{
			CustomerId: o.CustomerId,
			ItemId:     o.ItemId,
			CostEur:    o.CostEur,
		})
	}

	return items, nil
}

func (c *orderCache) GetSummaries() ([]Summary, error) {

	o, found := c.cache.Get(ordersSummaryCacheKey)

	if !found {
		return []Summary{}, errors.New("no summary found")
	}

	orderSummaries := o.(map[string]OrderSummaryCacheModel)

	var items []Summary
	for _, o := range orderSummaries {
		items = append(items, Summary{
			CustomerId:          o.CustomerId,
			NbrOfPurchasedItems: o.NbrOfPurchasedItems,
			TotalAmountEur:      o.TotalAmountEur,
		})
	}

	return items, nil
}
