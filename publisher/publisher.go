package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"
	"wb_l0/internal/app/models"

	"github.com/nats-io/stan.go"
)

func orderGenerator() *models.Order {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	id := randString(19, letters)

	order := models.Order{
		OrderUID:    id,
		TrackNumber: "WBL0",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+79998688999",
			Zip:     "263809",
			City:    "Moscow",
			Address: "Red Square 1",
			Region:  "Moscow",
			Email:   "test@yandex.ru",
		},
		Payment: models.Payment{
			Transaction:  "frhgtr",
			RequestID:    "fgb",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       9812,
			PaymentDt:    123455434,
			Bank:         "alpha",
			DeliveryCost: 1400,
			GoodsTotal:   343,
			CustomFee:    0,
		},
		Item: []models.Item{
			{
				ChrtID:      99345,
				TrackNumber: "WFEFWSFWE",
				Price:       234,
				Rid:         "wesfwefew",
				Name:        "sdvgdfrvg",
				Sale:        30,
				Size:        "0",
				TotalPrice:  345,
				NmID:        343545,
				Brand:       "Tst",
				Status:      202,
			},
		},
		Locale:            "ru",
		InternalSignature: "fsdg",
		CustomerID:        "test",
		DeliveryService:   "serr",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	return &order
}

func invalidOrderGenerator() *models.Order {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	id := randString(18, letters)

	order := models.Order{
		OrderUID:    id,
		TrackNumber: "WBL0",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+79998688999",
			Zip:     "263809",
			City:    "Moscow",
			Address: "Red Square 1",
			Region:  "Moscow",
			Email:   "test@yandex.ru",
		},
		Payment: models.Payment{
			Transaction:  "frhgtr",
			RequestID:    "fgb",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       9812,
			PaymentDt:    123455434,
			Bank:         "alpha",
			DeliveryCost: 1400,
			GoodsTotal:   343,
			CustomFee:    0,
		},
		Locale:            "ru",
		InternalSignature: "fsdg",
		CustomerID:        "test",
		DeliveryService:   "serr",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	return &order
}

func Publish() {
	sc, err := stan.Connect("test-cluster", "publisher")
	defer sc.Close()
	if err != nil {
		log.Printf("[Error]: publisher can't connect to Nats: %v\n", err)
		return
	}
	for {
		for i := 0; i < 5; i++ {
			order := orderGenerator()
			orderJs, err := json.Marshal(order)
			if err != nil {
				log.Printf("[Error]: publisher can't marshal to JSON: %v\n", err)
				return
			}

			err = sc.Publish("order", orderJs)
			if err != nil {
				log.Printf("[Error]: publisher can't publish: %v\n", err)
				return
			}
			log.Println("message send")
			time.Sleep(5 * time.Second)
		}

		inOrder := invalidOrderGenerator()
		inOrderJs, err := json.Marshal(inOrder)
		if err != nil {
			log.Printf("[Error]: publisher can't marshal to JSON: %v\n", err)
			return
		}

		err = sc.Publish("order", inOrderJs)
		if err != nil {
			log.Printf("[Error]: publisher can't publish: %v\n", err)
			return
		}
		log.Println("message send")
		time.Sleep(10 * time.Second)
	}
}

func main() {
	Publish()
}

func randString(n int, letters []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
