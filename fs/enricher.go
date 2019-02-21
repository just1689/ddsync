package fs

import (
	"github.com/fsnotify/fsnotify"
	"os"
)

type Enriched struct {
	InstanceID  string
	Event       *fsnotify.Event
	FullPath    string
	IsDirectory bool
	Directory   string
}

func (e *Enriched) Read() chan *Frame {
	return readSlowly(e.InstanceID, e.FullPath)
}

func StartEnrich(instanceID, directory string, in chan *fsnotify.Event) (out chan *Enriched) {
	out = make(chan *Enriched)
	go func() {
		for {
			i := <-in
			o := &Enriched{
				Event:      i,
				FullPath:   i.Name,
				InstanceID: instanceID,
			}
			o.IsDirectory = isDir(o.FullPath)
			out <- o
		}
	}()
	return
}

func isDir(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}
	return false
}
