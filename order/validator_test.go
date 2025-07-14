package order

import (
	"testing"

	"company.com/order-service/order/model"
	"github.com/stretchr/testify/assert"
)

var testRequest = model.CreateOrderRequest{
	Orders: []model.Order{
		{
			CustomerId: "1",
			OrderId:    "2",
			TimeStamp:  123,
			Items: []model.OrderItem{
				{
					ItemId:  "3",
					CostEur: 123.2,
				},
			},
		},
		{
			CustomerId: "2",
			OrderId:    "5",
			TimeStamp:  123,
			Items: []model.OrderItem{
				{
					ItemId:  "4",
					CostEur: 13.2,
				},
			},
		},
	},
}

func TestShouldReturnNoErrorWhenRequestIsValid(t *testing.T) {
	assertion := assert.New(t)

	err := validateOrderRequest(testRequest)

	assertion.ErrorIs(err, nil)
}
func TestShouldReturnErrorWhenCustomerIdIsEmpty(t *testing.T) {
	assertion := assert.New(t)

	testRequest.Orders[0].CustomerId = ""

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "customer id is missing in the request")
}

func TestShouldReturnErrorWhenOrderIdIsEmpty(t *testing.T) {
	assertion := assert.New(t)

	testRequest.Orders[0].OrderId = ""

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "order id is missing in the request")
}

func TestShouldReturnErrorWhenTimestampIsInvalid(t *testing.T) {
	assertion := assert.New(t)

	testRequest.Orders[0].TimeStamp = -1

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "invalid timestamp in the request")
}

func TestShouldReturnErrorWhenOrderItemIdIsEmpty(t *testing.T) {
	assertion := assert.New(t)

	testRequest.Orders[0].Items = append(testRequest.Orders[0].Items, model.OrderItem{
		ItemId:  "",
		CostEur: 123,
	})

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "item at index 1 does not have item id")
}

func TestShouldReturnErrorWhenOrderItemCostIsInvalid(t *testing.T) {
	assertion := assert.New(t)

	testRequest.Orders[0].Items = append(testRequest.Orders[0].Items, model.OrderItem{
		ItemId:  "1",
		CostEur: -12,
	})

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "item at index 1 has negative cost value")
}
