package cache

import (
	"errors"
	"testing"
	"time"

	"company.com/order-service/config"
	"company.com/order-service/internal/mocks"
	"company.com/order-service/order/model"

	goCache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func getOrders() []model.OrderCacheModel {
	return []model.OrderCacheModel{
		{
			CustomerId: "1",
			OrderId:    "2",
			ItemId:     "1",
			CostEur:    12.3,
			CreatedAt:  time.Now(),
		},
		{
			CustomerId: "1",
			OrderId:    "2",
			ItemId:     "1",
			CostEur:    12.3,
			CreatedAt:  time.Now(),
		},
	}
}

func getSummaries() map[string]model.OrderSummaryCacheModel {
	var summaries = make(map[string]model.OrderSummaryCacheModel)
	summaries["1"] = model.OrderSummaryCacheModel{
		CustomerId:          "1",
		TotalAmountEur:      24.6,
		NbrOfPurchasedItems: 2,
	}
	return summaries
}

func TestAddOrdersSuccessShouldSaveItemsInMemory(t *testing.T) {
	assertion := assert.New(t)
	orders := getOrders()
	summaries := getSummaries()

	defaultCacheExpiration := time.Duration(config.Config.CacheExpiryDuration) * time.Hour

	c := goCache.New(defaultCacheExpiration, time.Duration(config.Config.CacheCleanupInterval)*time.Hour)

	orderCache := NewOrderCache(c, defaultCacheExpiration)

	orderCache.AddOrders(orders, summaries)

	o, _ := c.Get(ordersCacheKey)

	result := o.([]model.OrderCacheModel)

	assertion.Equal(len(result), 2)
}

func TestAddOrdersFail(t *testing.T) {
	assertion := assert.New(t)
	orders := getOrders()
	summaries := getSummaries()

	defaultCacheExpiration := time.Duration(config.Config.CacheExpiryDuration) * time.Hour

	c := new(mocks.Cache)
	c.On("Get", mock.Anything).Return([]model.OrderCacheModel{}, false)
	c.On("Add", mock.Anything, mock.Anything, defaultCacheExpiration).Return(errors.New("err"))

	orderCache := NewOrderCache(c, defaultCacheExpiration)

	err := orderCache.AddOrders(orders, summaries)

	o, _ := c.Get(ordersCacheKey)

	result := o.([]model.OrderCacheModel)

	assertion.Equal(len(result), 0)

	assertion.Equal(err, errors.New("err"))

}

func TestGetOrdersSuccess(t *testing.T) {
	assertion := assert.New(t)
	orders := getOrders()

	defaultCacheExpiration := time.Duration(config.Config.CacheExpiryDuration) * time.Hour

	c := goCache.New(defaultCacheExpiration, time.Duration(config.Config.CacheCleanupInterval)*time.Hour)

	c.Add(ordersCacheKey, orders, defaultCacheExpiration)

	orderCache := NewOrderCache(c, defaultCacheExpiration)

	resp, err := orderCache.GetOrders()

	assertion.Equal(err, nil)

	assertion.Equal(len(resp), 2)
}

func TestGetOrdersFail(t *testing.T) {
	assertion := assert.New(t)

	defaultCacheExpiration := time.Duration(config.Config.CacheExpiryDuration) * time.Hour

	c := goCache.New(defaultCacheExpiration, time.Duration(config.Config.CacheCleanupInterval)*time.Hour)

	orderCache := NewOrderCache(c, defaultCacheExpiration)

	resp, err := orderCache.GetOrders()

	assertion.Equal(len(resp), 0)
	assertion.Equal(err, errors.New("no order found"))
}

func TestGetSummariesSuccess(t *testing.T) {
	assertion := assert.New(t)

	summaries := getSummaries()

	defaultCacheExpiration := time.Duration(config.Config.CacheExpiryDuration) * time.Hour

	c := goCache.New(defaultCacheExpiration, time.Duration(config.Config.CacheCleanupInterval)*time.Hour)

	c.Add(ordersSummaryCacheKey, summaries, defaultCacheExpiration)

	orderCache := NewOrderCache(c, defaultCacheExpiration)

	resp, err := orderCache.GetSummaries()

	assertion.Equal(err, nil)

	assertion.Equal(len(resp), 1)

	assertion.Equal(resp[0].NbrOfPurchasedItems, 2)
}

func TestGetSummariesFail(t *testing.T) {
	assertion := assert.New(t)

	defaultCacheExpiration := time.Duration(config.Config.CacheExpiryDuration) * time.Hour

	c := goCache.New(defaultCacheExpiration, time.Duration(config.Config.CacheCleanupInterval)*time.Hour)

	orderCache := NewOrderCache(c, defaultCacheExpiration)

	resp, err := orderCache.GetSummaries()

	assertion.Equal(len(resp), 0)
	assertion.Equal(err, errors.New("no summary found"))
}
