package io

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Frame struct {
	Buffer *[]byte
	Len    int
}

func readSlowly(filename string) (out chan *Frame) {
	out = make(chan *Frame)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			logrus.Error(err)
		}
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
			out <- frame
			fmt.Print(string(buffer))
		}
	}()
	return out
}
