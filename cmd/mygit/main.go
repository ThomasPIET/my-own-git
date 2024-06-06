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
	fmt.Println("Logs from your program will appear here!")

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	case "cat-file":
		if len(os.Args) < 3 {
			log.Panic("something is missing")
		}
		catFile()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

func catFile() {
	blob_sha := os.Args[3]
	path := fmt.Sprintf(".git/objects/%v/%v", blob_sha[0:2], blob_sha[2:])
	file, err := os.Open(path)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()
	zlib_r, err := zlib.NewReader(file)
	if err != nil {
		log.Panic(err)
	}
	read_zlib, err := io.ReadAll(zlib_r)
	if err != nil {
		log.Panic(err)
	}
	parts := strings.Split(string(read_zlib), "\x00")
	fmt.Print(parts[1])
}
