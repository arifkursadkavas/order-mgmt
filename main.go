package main

import (
	"context"
	"fmt"
	"time"

	"company.com/retail/config"
	"company.com/retail/order"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	if err := config.LoadConfig("./config"); err != nil {
		panic(fmt.Errorf("config error: %s", err))
	}

	r := gin.Default()

	rg := r.Group("/api/v1")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Config.DBConnectionString))

	if err != nil {
		panic(err)
	}

	orderDb := order.NewOrderStore(client)

	orderHandler := order.NewOrderHandler(rg, orderDb)

	orderHandler.RegisterRoutes(rg)

	r.Run(fmt.Sprintf(":%v", config.Config.ServerPort))

	r.Run()

}
