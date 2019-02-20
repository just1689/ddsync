package main

import (
	"flag"
	"fmt"
	"github.com/just1689/gsyncg/io"
	"strings"
)

var directories = flag.String("dirs", ".", "Directors separated by a comma.")

func main() {
	flag.Parse()
	done := make(chan bool)

	dirs := strings.Split(*directories, ",")

	for _, d := range dirs {
		events := io.Watch(".")
		enriched := io.StartEnrich(".", events)

		go func() {
			for e := range enriched {
				fmt.Println(e.FullPath, e.IsDirectory, e.Event.Name, e.Event.Op, e.Directory)
			}
		}()

	}
	<-done
}
