package order

type OrderItems struct {
	ItemId  string
	CostEur float32
}

type CreateOrderRequest struct {
	CustomerId string
	OrderId    string
	TimeStamp  string
	Items      []OrderItems
}
