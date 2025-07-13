package main

import (
	"fmt"
	"time"

	"company.com/order-service/config"
	"company.com/order-service/order"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func main() {

	if err := config.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("config error: %s", err))
	}

	r := gin.Default()
	rg := r.Group("/api/v1") //group routes with version v1

	defaultCacheExpiration := time.Duration(config.Config.CacheExpiryDuration) * time.Hour

	cache := cache.New(defaultCacheExpiration, time.Duration(config.Config.CacheCleanupInterval)*time.Hour)

	orderCache := order.NewOrderCache(cache, defaultCacheExpiration)

	orderHandler := order.NewOrderHandler(rg, orderCache)

	orderHandler.RegisterRoutes(rg)

	r.Run(fmt.Sprintf(":%v", config.Config.ServerPort)) // Starts server and waits for connection.
}
