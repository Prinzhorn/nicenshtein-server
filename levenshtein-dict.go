package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

const ASCII_LOWER_CASE_A_OFFSET = 97

type node struct {
	next [26]*node
	end  bool
}

var root [26]*node

func addToIndex(word string) {
	//fmt.Println(word)

	var currentNodeList *[26]*node = &root

	for characterIndex := 0; characterIndex < len(word); characterIndex++ {
		normalizedIndex := word[characterIndex] - ASCII_LOWER_CASE_A_OFFSET
		atEnd := characterIndex == len(word)-1

		//Add a new entry if the current character is not present.
		if currentNodeList[normalizedIndex] == nil {
			//fmt.Println(characterIndex, normalizedIndex)
			var newNodeList [26]*node
			currentNodeList[normalizedIndex] = &node{newNodeList, atEnd}
		}

		currentNodeList = &currentNodeList[normalizedIndex].next
	}
}

func findInIndex(word string) bool {
	fmt.Println("----------")
	fmt.Println(word)

	if len(word) == 0 {
		return false
	}

	var currentNodeList *[26]*node = &root
	var currentNode *node

	for characterIndex := 0; characterIndex < len(word); characterIndex++ {
		normalizedIndex := word[characterIndex] - ASCII_LOWER_CASE_A_OFFSET

		//Current character not found, it is not in the index.
		if currentNodeList[normalizedIndex] == nil {
			return false
		}

		currentNode = currentNodeList[normalizedIndex]
		currentNodeList = &currentNode.next
	}

	return currentNode.end
}

func main() {
	file, err := os.Open(os.Args[1])

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

	start := time.Now()

	indexed := findInIndex(os.Args[2])

	t := time.Now()
	elapsed := t.Sub(start)

	fmt.Println(indexed, elapsed)
}
