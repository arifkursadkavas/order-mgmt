package order

import "time"

// Request Model
type Order struct {
	ItemId  string  `json:"itemId"`
	CostEur float32 `json:"costEur"`
}

type CreateOrderRequest struct {
	CustomerId string  `json:"customerId"`
	OrderId    string  `json:"orderId"`
	TimeStamp  int64   `json:"timestamp"` // Note: I assumed this to be int64 instead of a string to eliminate a conversion on the backend. But can be for sure string
	Items      []Order `json:"items"`
}

// Response Models
type ListItemsResponse struct {
	Items []Item `json:"items"`
}
type Item struct {
	CustomerId string  `json:"customerId"`
	ItemId     string  `json:"itemId"`
	CostEur    float32 `json:"costEur"`
}

type ListSummariesResponse struct {
	Summaries []Summary `json:"summaries"`
}
type Summary struct {
	CustomerId          string  `json:"customerId"`
	NbrOfPurchasedItems int     `json:"nbrOfPurchasedItems"`
	TotalAmountEur      float32 `json:"totalAmountEur"`
}

// Storage Model
type OrderDb struct {
	CustomerId string    `bson:"customerId"`
	OrderId    string    `bson:"orderId"`
	ItemId     string    `bson:"itemId"`
	CostEur    float32   `bson:"costEur"`
	CreatedAt  time.Time `bson:"createdAt"`
}
