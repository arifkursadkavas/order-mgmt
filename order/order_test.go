package order

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"company.com/order-service/internal/mocks"
	"company.com/order-service/order/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Reviewcomments -- tests are not validating response body, they are just checking status codes

func getBody() model.CreateOrderRequest {
	return model.CreateOrderRequest{
		Orders: []model.Order{
			{
				CustomerId: "1",
				OrderId:    "2",
				TimeStamp:  1,
				Items: []model.OrderItem{
					{
						ItemId:  "1",
						CostEur: 12,
					},
				},
			},
		},
	}
}

func TestCreateOrders(t *testing.T) {
	assertion := assert.New(t)

	body := getBody()
	// Reviewcomments --- maybe put invalid body in a separate function
	invalidBody := getBody()
	invalidBody.Orders[0].CustomerId = ""

	tests := []struct {
		name  string
		body  model.CreateOrderRequest
		code  int
		cache *mocks.OrderCache
	}{
		{
			name: "TestCreateOrdersSuccess_200",
			body: body,
			code: 200,
			cache: func() *mocks.OrderCache {
				ch := new(mocks.OrderCache)
				ch.On("AddOrders", mock.Anything, mock.Anything).Return(nil)
				return ch
			}(),
		},
		{
			name: "TestCreateFailsWhenRequestBodyInvalid_400",
			body: model.CreateOrderRequest{},
			code: 400,
			cache: func() *mocks.OrderCache {
				ch := new(mocks.OrderCache)
				ch.On("AddOrders", mock.Anything, mock.Anything).Return(nil)
				return ch
			}(),
		},
		{
			name: "TestCreateFailsWhenRequestBodyContainsInvalidField_400",
			body: invalidBody,
			code: 400,
			cache: func() *mocks.OrderCache {
				ch := new(mocks.OrderCache)
				ch.On("AddOrders", mock.Anything, mock.Anything).Return(nil)
				return ch
			}(),
		},
		{
			name: "TestCreateFailsWhenRequestCacheStoreFails_500",
			body: body,
			code: 500,
			cache: func() *mocks.OrderCache {
				ch := new(mocks.OrderCache)
				ch.On("AddOrders", mock.Anything, mock.Anything).Return(errors.New("error"))
				return ch
			}(),
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()

		g := gin.New()
		h := NewOrderHandler(&g.RouterGroup, tt.cache)

		c, _ := gin.CreateTestContext(w)
		jsonbytes, err := json.Marshal(tt.body)
		if err != nil {
			panic(err)
		}

		c.Request = &http.Request{
			Method: http.MethodPost,
			Header: make(http.Header),
		}
		c.Request.Header.Set("Content-Type", "application/json")
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

		h.CreateOrders(c)
		assertion.Equal(w.Code, tt.code)
	}

}

func TestListItems(t *testing.T) {
	assertion := assert.New(t)

	tests := []struct {
		name  string
		code  int
		cache *mocks.OrderCache
	}{
		{
			name: "TestListItemsSuccess_200",

			code: 200,
			cache: func() *mocks.OrderCache {
				ch := new(mocks.OrderCache)
				ch.On("GetOrders").Return([]model.Item{}, nil)
				return ch
			}(),
		},
		{
			name: "TestListItemsFailWhenCacheFails500",
			code: 500,
			cache: func() *mocks.OrderCache {
				ch := new(mocks.OrderCache)
				ch.On("GetOrders").Return([]model.Item{}, errors.New("err"))
				return ch
			}(),
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()

		g := gin.New()
		h := NewOrderHandler(&g.RouterGroup, tt.cache)

		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Method: http.MethodGet,
			Header: make(http.Header),
		}
		c.Request.Header.Set("Content-Type", "application/json")

		h.ListItems(c)
		assertion.Equal(w.Code, tt.code)
	}
}

func TestListSummaries(t *testing.T) {
	assertion := assert.New(t)

	tests := []struct {
		name  string
		code  int
		cache *mocks.OrderCache
	}{
		{
			name: "TestListSummariesSuccess_200",

			code: 200,
			cache: func() *mocks.OrderCache {
				ch := new(mocks.OrderCache)
				ch.On("GetSummaries").Return([]model.Summary{}, nil)
				return ch
			}(),
		},
		{
			name: "TestListSummariesFailWhenCacheFails500",
			code: 500,
			cache: func() *mocks.OrderCache {
				ch := new(mocks.OrderCache)
				ch.On("GetSummaries").Return([]model.Summary{}, errors.New("err"))
				return ch
			}(),
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()

		g := gin.New()
		h := NewOrderHandler(&g.RouterGroup, tt.cache)

		c, _ := gin.CreateTestContext(w)

		c.Request = &http.Request{
			Method: http.MethodGet,
			Header: make(http.Header),
		}
		c.Request.Header.Set("Content-Type", "application/json")

		h.ListSummaries(c)
		assertion.Equal(w.Code, tt.code)
	}
}
