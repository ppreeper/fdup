package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var files = make(map[[sha1.Size]byte]string)

func checkDuplicate(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	if info.IsDir() {
		return nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Print(err)
		return nil
	}
	digest := sha1.Sum(data)
	if v, ok := files[digest]; ok {
		fmt.Printf("#rm %q\nrm %q\n\n", v, path)
	} else {
		files[digest] = path
	}

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
}
