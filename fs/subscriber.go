package fs

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

func CreateEventSubscriberForInstance(instanceID string) (f func(b []byte)) {
	f = func(b []byte) {
		i := instanceID
		var e = &Enriched{}
		err := json.Unmarshal(b, e)
		if err != nil {
			logrus.Error(err)
			//TODO: handle out of sync
			return
		}
		handleRemoteEnriched(i, e)
	}
	return

}

func handleRemoteEnriched(instanceID string, enriched *Enriched) {

	if enriched.Event.Op == fsnotify.Chmod {
		//Ignore CHMODs
		return
	} else if enriched.InstanceID == instanceID {
		//Ignore my own events
		return
	}

	fmt.Println("Received enriched:", enriched.Event.Op, enriched.FullPath)
	//TODO: make local changes
}
