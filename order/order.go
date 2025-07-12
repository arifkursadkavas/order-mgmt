package order

import (
	"context"
	"net/http"
	"time"

	"company.com/retail/config"
	"github.com/gin-gonic/gin"
)

type OrderI interface {
	createOrders(c *gin.Context)
	listItems(c *gin.Context)
	listSummaries(c *gin.Context)
	RegisterRoutes(r *gin.RouterGroup)
}

type OrderH struct {
	store OrderStore
}

func NewOrderHandler(r *gin.RouterGroup, store OrderStore) OrderI {

	return &OrderH{
		store: store,
	}

}

func (o *OrderH) createOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.APIDefaultTimeout)*time.Second)
	defer cancel()

	var request CreateOrderRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err = validateOrderRequest(request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var orders []interface{}

	for _, ord := range request.Items {
		orders = append(orders, OrderDb{
			CustomerId: request.CustomerId,
			OrderId:    request.OrderId,
			CreatedAt:  time.Unix(0, request.TimeStamp*int64(time.Millisecond)),
			ItemId:     ord.ItemId,
			CostEur:    ord.CostEur,
		})
	}

	err = o.store.AddOrders(ctx, orders)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (o *OrderH) listItems(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.APIDefaultTimeout)*time.Second)
	// defer cancel()

}

func (o *OrderH) listSummaries(c *gin.Context) {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.APIDefaultTimeout)*time.Second)
	// defer cancel()

}

func (o *OrderH) RegisterRoutes(r *gin.RouterGroup) {

	r.POST("/order", o.createOrders)
	r.GET("/item/list", o.listItems)
	r.GET("/summary/list", o.listSummaries)

}
