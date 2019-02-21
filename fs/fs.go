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

				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logrus.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(directory)
	if err != nil {
		logrus.Fatal(err)
	}
	return
}
