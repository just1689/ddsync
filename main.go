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
var ID string

const TopicEvent = "ddsync-event-dir"
const TopicFrame = "ddsync-frame-file"

func main() {

	setupID()

	flag.Parse()
	done := make(chan bool)
	dirs := strings.Split(*directories, ",")
	c := setupNSQ()

	for _, d := range dirs {

		events := fs.Watch(d)
		enriched := fs.StartEnrich(ID, d, events)

		go func() {
			for e := range enriched {
				b, err := json.Marshal(*e)
				if err != nil {
					logrus.Error(err)
					continue
				}
				err = c.Publish(TopicEvent, b)
				if err != nil {
					logrus.Error(err)
					continue
				}

				if !e.IsDirectory && e.Event.Op == fsnotify.Write {
					frames := e.Read()
					for f := range frames {
						b, err := json.Marshal(*f)
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

				logrus.Debugln(e.FullPath, e.IsDirectory, e.Event.Name, e.Event.Op, e.Directory)
			}
		}()

	}
	<-done
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
	c = nsq.Connect(*lookupAddress)
	err := c.AddHandler(TopicEvent, ID, fs.CreateEventSubscriberForInstance(ID))
	if err != nil {
		logrus.Error(err)
		return
	}
	return

}
