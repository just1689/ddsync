package fs

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Frame struct {
	Filename string
	Number   int
	Buffer   *[]byte
	Len      int
}

func readSlowly(filename string) (out chan *Frame) {
	logrus.Debugln("Reading", filename)
	out = make(chan *Frame)
	go func() {
		file, err := os.Open(filename)
		defer file.Close()
		if err != nil {
			logrus.Error(err)
		}
		logrus.Debugln("... starting", filename)
		n := 1
		for {
			buffer := make([]byte, 4096)
			l, err := file.Read(buffer)
			if err == io.EOF {
				break
			}
			frame := &Frame{
				Filename: filename,
				Number:   n,
				Buffer:   &buffer,
				Len:      l,
			}
			logrus.Debugln("... channeling", filename)
			out <- frame
			n += 1
		}
		close(out)
	}()
	return out
}
