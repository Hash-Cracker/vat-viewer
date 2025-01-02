package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if n > 0 {
			os.Stdout.Write(buffer[:n])
		}
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			}
			break
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		// No files provided; read from standard input
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		}
	} else {
		// Loop through all provided files
		for _, filename := range os.Args[1:] {
			printFile(filename)
		}
	}
}

