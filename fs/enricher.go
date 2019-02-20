package fs

import (
	"github.com/fsnotify/fsnotify"
	"os"
)

type Enriched struct {
	Event       *fsnotify.Event
	FullPath    string
	IsDirectory bool
	Directory   string
}

func (e *Enriched) Read() chan *Frame {
	return readSlowly(e.FullPath)
}

func StartEnrich(directory string, in chan *fsnotify.Event) (out chan *Enriched) {
	out = make(chan *Enriched)
	go func() {
		for {
			i := <-in
			o := &Enriched{
				Event: i,
				//FullPath: fmt.Sprintf("%s/%s", directory, i.Name), //UNIX vs WINDOWS?
				FullPath: i.Name,
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
