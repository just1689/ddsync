package fs

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"os"
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

func CreateFrameSubscriberForInstance(instanceID string) (f func(b []byte)) {
	f = func(b []byte) {
		i := instanceID
		var e = &Frame{}
		err := json.Unmarshal(b, e)
		if err != nil {
			logrus.Error(err)
			//TODO: handle out of sync
			return
		}
		handleRemoteFrame(i, e)
	}
	return

}

func handleRemoteFrame(instanceID string, frame *Frame) {

	if frame.InstanceID == instanceID {
		//Ignore my own events
		return
	}

	fmt.Println("Received frame:", frame.InstanceID, frame.Filename, frame.FrameNumber, frame.Len)

	if frame.FrameNumber == 1 {
		//TODO: delete the file and created it
	}
	f, err := os.OpenFile(frame.Filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		logrus.Error(err)
		//TODO: say Out of Sync
	}
	defer f.Close()
	//TODO: implement for where less than whole slice must be written. See `frame.Len`
	if _, err = f.Write(*frame.Buffer); err != nil {
		logrus.Error(err)
	}

}
