package apiserver

import (
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"sync"
	"wb_l0/internal/app/models"
	"wb_l0/internal/app/store"

	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	store *store.Store
	cache *store.Cache
}

func NewSubscriber(server *APIServer) *Subscriber {
	return &Subscriber{
		store: server.Store,
		cache: server.Store.Cache,
	}
}

func (s *Subscriber) Subscribe() {
	sc, err := stan.Connect("test-cluster", "subscriber")
	defer sc.Close()
	if err != nil {
		log.Printf("[Error]: subscriber can't connect to Nats: %v\n", err)
	}

	sc.Subscribe("order", func(msg *stan.Msg) {
		log.Println("message received")
		newOrder := models.Order{}
		ok := ValidateMessage(msg.Data)
		if !ok {
			return
		}
		err := json.Unmarshal(msg.Data, &newOrder)
		if err != nil {
			log.Println("[Error]: Can't unmarshal message")
			return
		}

		_, ok = s.cache.Get(newOrder.OrderUID)
		if ok {
			log.Println("This message already in cache")
			return
		}

		err = s.store.AddOrder(newOrder.OrderUID, msg.Data)
		if err != nil {
			log.Printf("[Error]: Can't add order: %v", err)
			return
		} else {
			log.Printf("Order added with id %v", newOrder.OrderUID)
		}
	})
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}

func ValidateMessage(m []byte) bool {

	var jsonData models.Order

	if err := json.Unmarshal(m, &jsonData); err != nil {
		return false
	}

	allFields := reflect.ValueOf(&jsonData).Elem()

	if allFields.NumField() != 14 {
		log.Printf("num field = %v", allFields.NumField())
		log.Println("[Error]: validation failed")
		return false
	}

	for i := 0; i < allFields.NumField(); i++ {
		if allFields.Field(i).IsZero() && strings.Contains(allFields.Type().Field(i).Tag.Get("validate"), "required") || len(jsonData.OrderUID) != 19 {
			log.Println("[Error]: validation failed")
			return false
		}
	}

	return true
}
