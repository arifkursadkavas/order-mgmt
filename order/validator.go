package order

import (
	"fmt"

	"company.com/order-service/order/model"
)

func validateOrderRequest(request model.CreateOrderRequest) error {

	for k, ord := range request.Orders {
		if ord.CustomerId == "" {
			return fmt.Errorf("customer id is missing in the request for order at index %d", k)
		}

		if ord.OrderId == "" {
			return fmt.Errorf("order id is missing in the request for order at index %d", k)
		}

		if ord.TimeStamp < 0 {
			return fmt.Errorf("invalid timestamp in the request for order at index %d", k)
		}

		if len(ord.Items) == 0 {
			return fmt.Errorf("no order item exist in the payload for order at index %d", k)
		}

		for i, item := range ord.Items {
			if item.ItemId == "" {
				return fmt.Errorf("item at index %d does not have item id for order at index %d", i, k)
			}

			if item.CostEur < 0 {
				return fmt.Errorf("item at index %d has negative cost value for order at index %d", i, k)
			}
		}

	}

	return nil
}
