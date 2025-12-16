package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"exc9/mapred"
)

// Main function
func main() {
	// todo read file
	file, err := os.Open("res/meditations.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var text []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// todo run your mapreduce algorithm
	var mr mapred.MapReduce
	results := mr.Run(text)

	// todo print your result to stdout
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}

	for _, word := range keys {
		fmt.Printf("%s: %d\n", word, results[word])
	}
}
