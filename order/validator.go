package order

import (
	"errors"
	"fmt"
)

func validateOrderRequest(request CreateOrderRequest) error {
	if request.CustomerId == "" {
		return errors.New("customer id is missing in the request")
	}

	if request.OrderId == "" {
		return errors.New("order id is missing in the request")
	}

	if request.TimeStamp < 0 {
		return errors.New("invalid timestamp in the request")
	}

	if len(request.Items) == 0 {
		return errors.New("no order item exist in the payload")
	}

	for i, item := range request.Items {
		if item.ItemId == "" {
			return fmt.Errorf("item at index %d does not have item id", i)
		}

		if item.CostEur < 0 {
			return fmt.Errorf("item at index %d has negative cost value", i)
		}
	}

	return nil
}
