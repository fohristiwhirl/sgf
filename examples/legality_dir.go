package main

// Scan a directory of SGF files for illegal moves.

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	sgf ".."
)

func main() {

	dirs := os.Args[1:]

	for _, d := range dirs {

		files, err := ioutil.ReadDir(d)

		if err != nil {
			panic(err.Error())
		}

		for _, f := range files {
			err := handle_file(d, f.Name())
			if err != nil {
				fmt.Printf("%s: %v\n", f.Name(), err)
			}
		}
	}
}

func handle_file(dirname, filename string) error {

	path := filepath.Join(dirname, filename)

	node, err := sgf.Load(path)
	if err != nil {
		return err
	}

	i := 0

	for {
		child := node.MainChild()
		if child == nil {
			break
		}

		i++

		board := node.Board()

		b, _ := child.GetValue("B")
		if b != "" && b != "tt" {
			_, err := board.LegalColour(b, sgf.BLACK)
			if err != nil {
				return fmt.Errorf("Move %d: %v", i, err)
			}
		}

		w, _ := child.GetValue("W")
		if w != "" && w != "tt" {
			_, err := board.LegalColour(w, sgf.WHITE)
			if err != nil {
				return fmt.Errorf("Move %d: %v", i, err)
			}
		}

		node = child
	}

	return nil
}
