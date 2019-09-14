package main

import (
	"fmt"
	"github.com/nats-io/go-nats"
	log "github.com/sirupsen/logrus"
	"os"
)

type BucketMessage struct {
	Bucket string
	Key    string
}

func init() {
	os.Setenv("NATS_URL", "nats://nats.app.sg.tb.you.co:4222")
	os.Setenv("FLUENTD_HOST", "")
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
	// Subscribe
	c := make(chan int)
	//t := make(chan *BucketMessage, 1)
	if _, err = nc.Subscribe("log.elb", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("acknowledged"))
	}); err != nil {
		log.Fatal(err)
	}
	if _, err = ec.Subscribe("log.elb", func(s *BucketMessage) {
		//t <- &BucketMessage{
		//	Bucket: s.Bucket,
		//	Key:    s.Key,
		//	}
		if err := LogParser(s); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s, %s", s.Bucket, s.Key)
	}); err != nil {
		log.Fatal(err)
	}
	//b := <- t
	<-c
}
