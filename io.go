package main

import (
	"image"
	"io"
	"os"
)

func readFiles(files []string, out chan<- image.Image, errChan chan<- error) {
	for _, file := range files {
		reader, err := os.Open(file)
		if err != nil {
			errChan <- err
			continue
		}
		defer reader.Close()

		img, _, err := image.Decode(reader)
		if err != nil {
			errChan <- err
			continue
		}
		out <- img
	}
}

func readStdin(out chan<- image.Image, errChan chan<- error) {
	for {
		img, _, err := image.Decode(os.Stdin)
		if err != nil {
			if err == io.EOF {
				break
			}

			errChan <- err
			continue
		}

		out <- img
	}
}
