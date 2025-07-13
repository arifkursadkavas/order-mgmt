package order

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"company.com/order-service/config"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.APIDefaultTimeout)*time.Second)
	defer cancel()

	pageSize, pageNumber := getPaginationParams(c)

	items, err := o.store.GetOrders(ctx, pageSize, pageNumber)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ListItemsResponse{Items: items})

}

func (o *OrderH) listSummaries(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.APIDefaultTimeout)*time.Second)
	defer cancel()

	pageSize, pageNumber := getPaginationParams(c)

	summaries, err := o.store.GetSummaries(ctx, pageSize, pageNumber)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, ListSummariesResponse{Summaries: summaries})

}

func (o *OrderH) RegisterRoutes(r *gin.RouterGroup) {

	r.POST("/order", o.createOrders)
	r.GET("/item/list", o.listItems)
	r.GET("/summary/list", o.listSummaries)

}

// Gets pagination parameters,
// Tries to obtain from query param of the request,
// Falls back to default values if not provided.
// returns pageSize, pageNumber respectively.
func getPaginationParams(c *gin.Context) (int, int) {

	pageSize := 10  //default page size
	pageNumber := 1 //default page number, first page

	queryPageSize := c.Query("pageSize")
	queryPageNumber := c.Query("pageNumber")

	pageSize, err := strconv.Atoi(queryPageSize)
	if err != nil {
		log.Printf("page size is not provided, using default page size 10")
	}

	pageNumber, err = strconv.Atoi(queryPageNumber)
	if err != nil {
		log.Printf("page number is not provided, using default page number 1")
	}

	return pageSize, pageNumber
}
