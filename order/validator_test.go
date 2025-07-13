package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testRequest = CreateOrderRequest{
	CustomerId: "1",
	OrderId:    "2",
	TimeStamp:  123,
	Items: []Order{
		{
			ItemId:  "3",
			CostEur: 123.2,
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

	testRequest.CustomerId = ""

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "customer id is missing in the request")
}

func TestShouldReturnErrorWhenOrderIdIsEmpty(t *testing.T) {
	assertion := assert.New(t)

	testRequest.OrderId = ""

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "order id is missing in the request")
}

func TestShouldReturnErrorWhenTimestampIsInvalid(t *testing.T) {
	assertion := assert.New(t)

	testRequest.TimeStamp = -1

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "invalid timestamp in the request")
}

func TestShouldReturnErrorWhenOrderItemIdIsEmpty(t *testing.T) {
	assertion := assert.New(t)

	testRequest.Items = append(testRequest.Items, Order{
		ItemId:  "",
		CostEur: 123,
	})

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "item at index 1 does not have item id")
}

func TestShouldReturnErrorWhenOrderItemCostIsInvalid(t *testing.T) {
	assertion := assert.New(t)

	testRequest.Items = append(testRequest.Items, Order{
		ItemId:  "1",
		CostEur: -12,
	})

	err := validateOrderRequest(testRequest)

	assertion.ErrorContains(err, "item at index 1 has negative cost value")
}
