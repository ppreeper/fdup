package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var files = make(map[[sha1.Size]byte][]string)

func checkDuplicate(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	if info.IsDir() {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Print(err)
		return nil
	}
	digest := sha1.Sum(data)
	files[digest] = append(files[digest], path)

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	var dir string
	if len(os.Args) == 1 {
		dir = "."
	} else {
		dir = os.Args[1]
	}
	err := filepath.Walk(dir, checkDuplicate)
	if err != nil {
		log.Fatal(err)
	}
	resfiles := make(map[[sha1.Size]byte][]string)
	for digest, v := range files {
		if len(v) > 1 {
			resfiles[digest] = v
		}
	}
	for _, filelist := range resfiles {
		for _, filename := range filelist {
			fmt.Println("./" + filename)
		}
		fmt.Println()
	}
}
