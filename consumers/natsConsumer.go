package consumers

import (
	"awesomeProject/cache"
	"awesomeProject/repository"
	"awesomeProject/validators"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"os"
	"os/signal"
	"sync"
)

func NatsConsumer(wg *sync.WaitGroup, orderRepository *repository.OrderRepository, cacheStore *cache.InMemoryStore) {
	defer wg.Done()

	clusterID := os.Getenv("CLUSTER_ID")
	clientID := os.Getenv("CLIENT_CONSUMER_ID")
	natsURL := os.Getenv("NATS_URL")

	// Подключение к серверу NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal(err)
	}

	// Подключение к NATS Streaming
	sc, err := stan.Connect(clusterID, clientID, stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err)
	}

	// Обработка сигнала завершения работы
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nПолучен сигнал завершения. Завершение работы.")
		sc.Close()
		nc.Close()
		os.Exit(0)
	}()

	// Подписка на канал
	var _ stan.Subscription

	_, err = sc.Subscribe("test-channel", func(msg *stan.Msg) {

		order, err := validators.MessageValidator(string(msg.Data))
		if err != nil {
			fmt.Printf("Ошибка валидации: %v\n", err)
			return
		}

		err = orderRepository.CreateOrder(order)
		if err != nil {
			return
		}

		cacheStore.AddOrder(*order)

		fmt.Printf("Получено сообщение: %s\n", *order)

	}, stan.DurableName("test-durable"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Ожидание сообщений. Для завершения нажмите Ctrl+C.")
	select {}
}
