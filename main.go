package main

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/just1689/ddsync/fs"
	"github.com/just1689/ddsync/nsq"
	"github.com/sirupsen/logrus"
	"strings"
)

var directories = flag.String("dirs", ".", "Directors separated by a comma.")

func main() {
	flag.Parse()
	done := make(chan bool)

	dirs := strings.Split(*directories, ",")

	f := func(b []byte) {
		fmt.Println(">", string(b))
	}

	c := nsq.Connect("192.168.88.24:4160")
	err := c.AddHandler("a", "a", f)
	if err != nil {
		logrus.Error(err)
		return
	}
	err = c.Publish("a", []byte("something"))
	if err != nil {
		logrus.Error(err)
		return
	}

	for _, d := range dirs {

		events := fs.Watch(d)
		enriched := fs.StartEnrich(d, events)

		go func() {
			for e := range enriched {
				if !e.IsDirectory && e.Event.Op == fsnotify.Write {
					c := e.Read()
					for f := range c {
						logrus.Print(string(*f.Buffer))
					}

				}

				logrus.Println(e.FullPath, e.IsDirectory, e.Event.Name, e.Event.Op, e.Directory)
			}
		}()

	}
	<-done
}
