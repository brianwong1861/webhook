package main

import (
	"github.com/labstack/gommon/log"
	"github.com/nats-io/go-nats"
)
const  (
	NATS_URL = "nats://nats.app.sg.tb.you.co:4222"
)
func main(){
	nc, err := nats.Connect(NATS_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}
	defer ec.Close()

	// Define object
	type BucketMessage struct {
		Bucket string
		Key  string
	}
	// Subscribe
	c := make(chan bool)
	if _, err = nc.Subscribe("log.elb", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("Received"))
	} );err != nil {
		log.Fatal(err)
	}
	if _, err = ec.Subscribe("log.elb", func(s *BucketMessage) {
		log.Printf("Stock: %s - Price: %v", s.Bucket, s.Key)
	});err != nil {
		log.Fatal(err)
	}
	<- c
}
