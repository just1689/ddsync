package main

import (
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/just1689/ddsync/io"
	"github.com/sirupsen/logrus"
	"strings"
)

var directories = flag.String("dirs", ".", "Directors separated by a comma.")

func main() {
	flag.Parse()
	done := make(chan bool)

	dirs := strings.Split(*directories, ",")

	for _, d := range dirs {

		events := io.Watch(d)
		enriched := io.StartEnrich(d, events)

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
