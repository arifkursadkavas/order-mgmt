package order

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"company.com/order-service/config"
	"company.com/order-service/order/model"
	"github.com/gin-gonic/gin"
)

type OrderI interface {
	CreateOrders(c *gin.Context)
	ListItems(c *gin.Context)
	ListSummaries(c *gin.Context)
	RegisterRoutes(r *gin.RouterGroup)
}

type OrderH struct {
	cache OrderCache
}

func NewOrderHandler(r *gin.RouterGroup, cache OrderCache) OrderI {
	return &OrderH{
		cache: cache,
	}
}

func (o *OrderH) CreateOrders(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.APIDefaultTimeout)*time.Second)
	defer cancel()

	var request model.CreateOrderRequest

	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if len(request.Orders) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("request contains no order"))
		return
	}

	err = validateOrderRequest(request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	orders, orderSummaries := prepareOrderData(request.Orders)

	err = o.cache.AddOrders(orders, orderSummaries)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (o *OrderH) ListItems(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.APIDefaultTimeout)*time.Second)
	defer cancel()

	items, err := o.cache.GetOrders()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, model.ListItemsResponse{Items: items})

}

func (o *OrderH) ListSummaries(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), time.Duration(config.Config.APIDefaultTimeout)*time.Second)
	defer cancel()

	summaries, err := o.cache.GetSummaries()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, model.ListSummariesResponse{Summaries: summaries})

}

func (o *OrderH) RegisterRoutes(r *gin.RouterGroup) {

	r.POST("/order", o.CreateOrders)
	r.GET("/item/list", o.ListItems)
	r.GET("/summary/list", o.ListSummaries)

}

func prepareOrderData(ords []model.Order) ([]model.OrderCacheModel, map[string]model.OrderSummaryCacheModel) {

	var orders []model.OrderCacheModel
	var orderSummary = make(map[string]model.OrderSummaryCacheModel)

	for _, ord := range ords {

		currentOrderSumEur := float32(0.0) //To keep sum of all costs and use it in summary.

		for _, item := range ord.Items {
			orders = append(orders, model.OrderCacheModel{
				CustomerId: ord.CustomerId,
				OrderId:    ord.OrderId,
				CreatedAt:  time.Unix(0, ord.TimeStamp*int64(time.Millisecond)),
				ItemId:     item.ItemId,
				CostEur:    item.CostEur,
			})
			currentOrderSumEur += item.CostEur
		}

		//Calculate the summaries by grouping them under same customerId in case a request contains more than
		// one record(OrderCacheModel) for a customer.
		if os, found := orderSummary[ord.CustomerId]; found {
			updatedOs := model.OrderSummaryCacheModel{}
			updatedOs.NbrOfPurchasedItems = os.NbrOfPurchasedItems + len(ord.Items)
			updatedOs.TotalAmountEur = os.TotalAmountEur + currentOrderSumEur
			orderSummary[ord.CustomerId] = updatedOs
		} else {
			newOs := model.OrderSummaryCacheModel{
				NbrOfPurchasedItems: len(ord.Items),
				TotalAmountEur:      currentOrderSumEur,
			}
			orderSummary[ord.CustomerId] = newOs
		}
	}

	return orders, orderSummary

}
