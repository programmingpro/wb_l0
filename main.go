package main

import (
	"awesomeProject/cache"
	"awesomeProject/consumers"
	"awesomeProject/models"
	"awesomeProject/producers"
	"awesomeProject/repository"
	"awesomeProject/web"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	loadEnv()

	var wg sync.WaitGroup

	dsn := os.Getenv("DBSTRING")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	orderRepository := repository.NewOrderRepository(db)
	cacheStore := cache.NewInMemoryStore()

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.Order{}, &models.Delivery{}, &models.Payment{}, &models.Item{})

	wg.Add(3)

	orders, err := orderRepository.GetAllOrders()
	if err != nil {
		log.Fatal(err)
	}

	for _, order := range orders {
		cacheStore.AddOrder(order)
	}

	fmt.Println("Cache Data")
	fmt.Println(cacheStore.GetAllOrders())

	webServer := &web.WebServer{
		&cacheStore,
	}

	go webServer.Run(&wg)
	go consumers.NatsConsumer(&wg, &orderRepository, &cacheStore)
	go producers.NatsProducer()

	wg.Wait()
}
