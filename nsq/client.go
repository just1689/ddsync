package nsq

import (
	"github.com/nsqio/go-nsq"
	"github.com/nsqio/nsq/nsqd"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type NsqClient struct {
	lookupAddress string
	stop          chan bool
	p             *nsq.Producer
	cfg           *nsq.Config
}

func (n *NsqClient) Stop() {
	n.stop <- true
}

func (n *NsqClient) Publish(topic string, msg []byte) (err error) {

	err = n.p.Publish(topic, msg)
	if err != nil {
		logrus.Error(err)
	}
	return
}

func (n *NsqClient) AddHandler(topic, channel string, h func(b []byte)) (err error) {
	c, err := nsq.NewConsumer(topic, channel, n.cfg)
	if err != nil {
		log.Fatal(err)
	}
	c.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		h(m.Body)
		return nil
	}))
	err = c.ConnectToNSQD("localhost:4150")
	return
}

func Connect(lo string) (nsqClient *NsqClient) {

	nsqClient = &NsqClient{
		stop:          make(chan bool),
		lookupAddress: lo,
		cfg:           nsq.NewConfig(),
	}

	go func() {

		opts := nsqd.NewOptions()
		opts.NSQLookupdTCPAddresses = []string{lo}
		opts.MemQueueSize = int64(256) // 256 X 4KB = 1MB
		daemon := nsqd.New(opts)
		daemon.Main()
		<-nsqClient.stop
		daemon.Exit()
	}()

	time.Sleep(1 * time.Second)

	// Set up a Producer, pointing at the default host:port
	var err error
	nsqClient.p, err = nsq.NewProducer("localhost:4150", nsqClient.cfg)
	if err != nil {
		logrus.Error(err)
	}
	return

}
