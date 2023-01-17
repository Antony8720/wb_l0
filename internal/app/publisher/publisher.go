package publisher

import (
	"encoding/json"
	"log"
	"time"
	"wb_l0/internal/app/models"

	"github.com/google/uuid"
	"github.com/nats-io/stan.go"
)

func orderGenerator() *models.DataJSON {
	id := uuid.New()

	order := models.DataJSON{
		OrderUID:    id.String(),
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


func Publish() {
	sc, err := stan.Connect("wb_l0", "publisher")
	defer sc.Close()
	if err != nil {
		log.Printf("[Error]: publisher can't connect to Nats: %v\n", err)
		return
	}
	
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
}