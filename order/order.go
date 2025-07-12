package order

import "github.com/gin-gonic/gin"

type Order struct {
	DBClient *Client
}

func NewOrderHandler(r *gin.RouterGroup, db *Client) *Order {

	return &Order{
		DBClient: db,
	}

}

func (o *Order) createOrder(c *gin.Context) {

}

func (o *Order) RegisterRoutes(r *gin.RouterGroup) {

	r.POST("/order", o.createOrder)

}
