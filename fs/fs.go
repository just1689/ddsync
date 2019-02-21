package fs

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

func Watch(directory string) (events chan *fsnotify.Event) {
	events = make(chan *fsnotify.Event)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	go func() {
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				events <- &event
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Error(err)
			}
		}
	}()
	err = watcher.Add(directory)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}
