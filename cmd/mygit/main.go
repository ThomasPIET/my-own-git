package main

import (
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		_, err := fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		if err != nil {
			log.Panic(err)
		}
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				_, err := fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
				if err != nil {
					log.Panic(err)
				}
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			if _, fprintfErr := fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err); fprintfErr != nil {
				log.Panic(fprintfErr)
			}
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		if len(os.Args) < 3 {
			log.Panic("something is missing")
		}
		catFile()

	default:
		_, err := fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		if err != nil {
			log.Panic(err)
		}
		os.Exit(1)
	}
}

func catFile() {
	blobSha := os.Args[3]
	path := fmt.Sprintf(".git/objects/%v/%v", blobSha[0:2], blobSha[2:])
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Panic(err)
		}
	}(file)
	zlibR, err := zlib.NewReader(file)
	if err != nil {
		log.Panic(err)
	}
	readZlib, err := io.ReadAll(zlibR)
	if err != nil {
		log.Panic(err)
	}
	parts := strings.Split(string(readZlib), "\x00")
	fmt.Print(parts[1])
}
