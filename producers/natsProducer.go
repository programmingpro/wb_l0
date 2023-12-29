package producers

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

func NatsProducer() {

	clusterID := os.Getenv("CLUSTER_ID")
	clientID := os.Getenv("CLIENT_PRODUCER_ID")
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

	filePath := "messages/test.json"

	// Отправка сообщения в канал
	go func() {
		// Считывание данных из файла
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Ошибка при чтении файла %s: %v", filePath, err)
			return
		}

		// Декодирование JSON
		var messages []map[string]interface{}
		if err := json.Unmarshal(fileContent, &messages); err != nil {
			log.Fatalf("Ошибка при декодировании JSON из файла %s: %v", filePath, err)
			return
		}

		// Отправка сообщений в канал
		for _, message := range messages {
			jsonMessage, err := json.Marshal(message)
			if err != nil {
				log.Printf("Ошибка при маршалинге сообщения: %v\n", err)
				continue
			}

			err = sc.Publish("test-channel", []byte(jsonMessage))
			if err != nil {
				log.Printf("Ошибка при отправке сообщения: %v\n", err)
			} else {
				fmt.Printf("Отправлено сообщение\n")
			}

			// Пауза между отправками сообщений (можно регулировать)
			time.Sleep(2 * time.Second)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\nПолучен сигнал завершения. Завершение работы.")
		sc.Close()
		nc.Close()
		os.Exit(0)
	}()
}
