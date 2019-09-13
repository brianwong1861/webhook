package main

import (
	"github.com/labstack/gommon/log"
	"github.com/nats-io/go-nats"
	"os"
)

type BucketMessage struct {
	Bucket string
	Key    string
}

func main() {
	nc, err := nats.Connect(os.Getenv("NATS_URL")) // nats://localhost:4222
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
	// Subscribe
	c := make(chan bool)
	if _, err = nc.Subscribe("log.elb", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("Received"))
	}); err != nil {
		log.Fatal(err)
	}
	if _, err = ec.Subscribe("log.elb", func(s *BucketMessage) {
		LogParser(&BucketMessage{
			Bucket: s.Bucket,
			Key:    s.Key,
		})
	}); err != nil {
		log.Fatal(err)
	}
	<-c
}
