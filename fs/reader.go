package fs

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Frame struct {
	Buffer *[]byte
	Len    int
}

func readSlowly(filename string) (out chan *Frame) {
	logrus.Debugln("Reading", filename)
	out = make(chan *Frame)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Debugln("... starting", filename)
		for {
			buffer := make([]byte, 4096)
			l, err := file.Read(buffer)
			if err == io.EOF {
				break
			}
			frame := &Frame{
				Buffer: &buffer,
				Len:    l,
			}
			logrus.Debugln("... channeling", filename)
			out <- frame
		}
		file.Close()
		close(out)
	}()
	return out
}
