package subscriber

import (
	"log"

	"github.com/nats-io/stan.go"
)

func Subscriber () {
	sc, err := stan.Connect("wb_l0", "subscriber")
	defer sc.Close()
	if err != nil {
		log.Printf("[Error]: subscriber can't connect to Nats: %v\n", err)
	}

	sub, err := sc.Subscribe("order", func(msg *stan.Msg) {
		
	})
}