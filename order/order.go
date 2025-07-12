package order

import (
	"net/http"

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
	var request CreateOrderRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	err = validateOrderRequest(request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	o.store.AddOrders()
}

func (o *OrderH) listItems(c *gin.Context) {

}

func (o *OrderH) listSummaries(c *gin.Context) {

}

func (o *OrderH) RegisterRoutes(r *gin.RouterGroup) {

	r.POST("/order", o.createOrders)
	r.GET("/item/list", o.listItems)
	r.GET("/summary/list", o.listSummaries)

}
