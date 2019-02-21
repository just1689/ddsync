package fs

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Frame struct {
	InstanceID  string
	Filename    string
	FrameNumber int
	Buffer      *[]byte
	Len         int
}

func readSlowly(instanceID, filename string) (out chan *Frame) {
	logrus.Debugln("Reading", filename)
	out = make(chan *Frame)
	go func() {
		file, err := os.Open(filename)
		defer file.Close()
		if err != nil {
			logrus.Error(err)
		}
		logrus.Debugln("... starting", filename)
		frameNumber := 1
		for {
			buffer := make([]byte, 4096)
			l, err := file.Read(buffer)
			if err == io.EOF {
				break
			}
			frame := &Frame{
				InstanceID:  instanceID,
				Filename:    filename,
				FrameNumber: frameNumber,
				Buffer:      &buffer,
				Len:         l,
			}
			logrus.Debugln("... channeling", filename)
			out <- frame
			frameNumber += 1
		}
		close(out)
	}()
	return out
}
