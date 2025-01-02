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
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			break
		}
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		}
	} else {
		for _, filename := range args[1:] {
			printFile(filename)
		}
	}
}

