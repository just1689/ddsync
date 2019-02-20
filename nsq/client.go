package nsq

import (
	"github.com/nsqio/go-nsq"
	"github.com/nsqio/nsq/nsqd"
	"github.com/sirupsen/logrus"
	"log"
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
	err = n.p.Publish("embedded", msg)
	if err != nil {
		logrus.Error(err)
	}
	return
}

func (n *NsqClient) AddHandler(topic, channel string, h func(b []byte)) (err error) {
	// Now set up a consumer
	c, err := nsq.NewConsumer("embedded", "local", n.cfg)
	if err != nil {
		log.Fatal(err)
	}
	c.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		h(m.Body)
		return nil
	}))
	err = c.ConnectToNSQD(n.lookupAddress)
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
		daemon := nsqd.New(opts)
		daemon.Main()
		<-nsqClient.stop
		daemon.Exit()
	}()

	return

}
