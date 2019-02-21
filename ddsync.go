package main

import (
	"encoding/json"
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/google/uuid"
	"github.com/just1689/ddsync/fs"
	"github.com/just1689/ddsync/nsq"
	"github.com/sirupsen/logrus"
	"strings"
)

var directories = flag.String("dirs", ".", "Directors separated by a comma.")
var lookupAddress = flag.String("lookup", "", "Lookup address and port host:4160")
var listenAddress = flag.String("listen", "localhost:4150", "Address to host NSQ Daemon server on.")
var ID string

const TopicEvent = "ddsync-event-dir"
const TopicFrame = "ddsync-frame-file"

func main() {

	flag.Parse()
	setupID()
	dirs := strings.Split(*directories, ",")

	done := make(chan bool)
	c := setupNSQ()
	MonitorDirs(dirs, c)
	<-done

}

func MonitorDirs(dirs []string, c *nsq.NsqClient) {

	for _, d := range dirs {
		events := fs.Watch(d)
		enriched := fs.StartEnrich(ID, d, events)
		go func() {
			for e := range enriched {

				if err := publishEvent(c, e); err != nil {
					continue
				}

				if err := publishFrame(c, e); err != nil {
					continue
				}

				logrus.Debugln(e.FullPath, e.IsDirectory, e.Event.Name, e.Event.Op, e.Directory)
			}
		}()
	}
}

func publishFrame(c *nsq.NsqClient, e *fs.Enriched) (err error) {
	if !e.IsDirectory && e.Event.Op == fsnotify.Write {
		frames := e.Read()
		var b []byte
		for f := range frames {
			b, err = json.Marshal(*f)
			if err != nil {
				logrus.Error(err)
				continue
			}
			err = c.Publish(TopicFrame, b)
			if err != nil {
				logrus.Error(err)
				return
			}
			//logrus.Print(string(*f.Buffer))
		}

	}
	return
}

func publishEvent(c *nsq.NsqClient, e *fs.Enriched) (err error) {
	b, err := json.Marshal(*e)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = c.Publish(TopicEvent, b)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}

func setupID() {
	myID, err := uuid.NewRandom()
	if err != nil {
		logrus.Fatal(err)
	}
	ID = myID.String()
}

func setupNSQ() (c *nsq.NsqClient) {
	// Setup NSQ
	c = nsq.StartNSQDaemon(*lookupAddress, *listenAddress)
	err := c.AddHandler(TopicEvent, ID, fs.CreateEventSubscriberForInstance(ID))
	if err != nil {
		logrus.Error(err)
		panic(err)
		return
	}
	err = c.AddHandler(TopicFrame, ID, fs.CreateFrameSubscriberForInstance(ID))
	if err != nil {
		logrus.Error(err)
		panic(err)
		return
	}

	return

}
