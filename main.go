package main

import (
	"fmt"
	"time"

	"company.com/order-service/config"
	"company.com/order-service/order"
	cache "company.com/order-service/order/cache"
	"github.com/gin-gonic/gin"
	goCache "github.com/patrickmn/go-cache"
)

func main() {

	if err := config.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("config error: %s", err))
	}

	r := gin.Default()
	rg := r.Group("/api/v1") //group routes with version v1

	defaultCacheExpiration := time.Duration(config.Config.CacheExpiryDuration) * time.Hour

	c := goCache.New(defaultCacheExpiration, time.Duration(config.Config.CacheCleanupInterval)*time.Hour)

	orderCache := cache.NewOrderCache(c, defaultCacheExpiration)

	orderHandler := order.NewOrderHandler(rg, orderCache)

	orderHandler.RegisterRoutes(rg)

	r.Run(fmt.Sprintf(":%v", config.Config.ServerPort)) // Starts server and waits for connection.
}
