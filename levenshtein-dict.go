package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"unicode/utf8"
)

type node struct {
	children map[rune]*node
	word     string
}

var root node = node{make(map[rune]*node), ""}

func indexFile(fileName string) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		addToIndex(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func addToIndex(word string) {
	var currentNode *node = &root

	for index, runeValue := range word {
		childNode, ok := currentNode.children[runeValue]

		//We have not indexed this rune yet, create a new entry.
		if !ok {
			childNode = &node{make(map[rune]*node), ""}
			currentNode.children[runeValue] = childNode
		}

		//The node at the end of a word stores the full word, which also marks the end.
		//This makes the index less memory efficient, but vastly improves query performance.
		//Otherwise each query would need to collect the runes along the path and concat the word.
		if index == len(word)-len(string(runeValue)) {
			childNode.word = word
		}

		currentNode = childNode
	}
}

func findInIndex(word string) bool {
	var currentNode *node = &root

	for _, runeValue := range word {
		childNode, ok := currentNode.children[runeValue]

		//Current rune not indexed.
		if !ok {
			return false
		}

		currentNode = childNode
	}

	//Does a string terminate at this node?
	return currentNode.word != ""
}

func collectFromIndex(out *map[string]byte, currentNode *node, word string, distance byte, maxDistance byte) {
	if len(word) == 0 {
		if currentNode.word != "" {
			value, ok := (*out)[currentNode.word]

			//We have not seen this word or we have found a smaller distance.
			if !ok || distance < value {
				(*out)[currentNode.word] = distance
			}
		}

		return
	}

	if distance < maxDistance {
		for runeValue, _ := range currentNode.children {
			//Substitution (replace the first character with the current one).
			collectFromIndex(out, currentNode, string(runeValue)+word[1:], distance+1, maxDistance)

			//Insertion (add the current character as prefix).
			collectFromIndex(out, currentNode, string(runeValue)+word, distance+1, maxDistance)
		}

		//Deletion (skip first character).
		collectFromIndex(out, currentNode, word[1:], distance+1, maxDistance)
	}

	runeValue, _ := utf8.DecodeRuneInString(word[0:])
	nextNode := currentNode.children[runeValue]

	if nextNode != nil {
		//Move forward by one character without incrementing the distance.
		collectFromIndex(out, nextNode, word[1:], distance, maxDistance)
	}
}

func main() {
	start := time.Now()
	indexFile(os.Args[1])
	elapsed := time.Now().Sub(start)
	fmt.Println(elapsed)

	maxDistance, err := strconv.ParseInt(os.Args[3], 10, 8)

	if err != nil {
		fmt.Println(err)
	}

	out := make(map[string]byte)

	start = time.Now()

	collectFromIndex(&out, &root, os.Args[2], 0, byte(maxDistance))

	elapsed = time.Now().Sub(start)

	for key, value := range out {
		fmt.Println(key, value)
	}

	fmt.Println(elapsed)
}
