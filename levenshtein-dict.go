package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const ASCII_LOWER_CASE_A_OFFSET = 97

type node struct {
	next [26]*node
	word string
}

var root node

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
	//TODO: use a big ass slice of the root instead, like 26 pieces.
	var currentNodeList *[26]*node = &root.next

	for characterIndex := 0; characterIndex < len(word); characterIndex++ {
		normalizedIndex := word[characterIndex] - ASCII_LOWER_CASE_A_OFFSET

		//Add a new entry if the current character is not present.
		if currentNodeList[normalizedIndex] == nil {
			var newNodeList [26]*node

			//The node at the end of the word stores the full word.
			//This makes the index less memory efficient, but vastly improves query performance.
			if characterIndex == len(word)-1 {
				currentNodeList[normalizedIndex] = &node{newNodeList, word}
			} else {
				currentNodeList[normalizedIndex] = &node{newNodeList, ""}
			}
		}

		currentNodeList = &currentNodeList[normalizedIndex].next
	}
}

func findInIndex(word string) bool {
	if len(word) == 0 {
		return false
	}

	var currentNodeList *[26]*node = &root.next
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
		for i := 0; i < 26; i++ {
			iNode := currentNode.next[i]

			//There is no node for this character (no word in the index has this prefix).
			if iNode == nil {
				continue
			}

			currentCharacter := string(i + ASCII_LOWER_CASE_A_OFFSET)

			//Substitution (replace the first character with the current one).
			collectFromIndex(out, currentNode, currentCharacter+word[1:], distance+1, maxDistance)

			//Insertion (add the current character as prefix).
			collectFromIndex(out, currentNode, currentCharacter+word, distance+1, maxDistance)
		}

		//Deletion (skip first character).
		collectFromIndex(out, currentNode, word[1:], distance+1, maxDistance)
	}

	//Move forward by one character without incrementing the distance.
	normalizedIndex := word[0] - ASCII_LOWER_CASE_A_OFFSET
	nextNode := currentNode.next[normalizedIndex]

	if nextNode != nil {
		collectFromIndex(out, nextNode, word[1:], distance, maxDistance)
	}
}

func main() {
	indexFile(os.Args[1])

	maxDistance, err := strconv.ParseInt(os.Args[3], 10, 8)

	if err != nil {
		fmt.Println(err)
	}

	out := make(map[string]byte)

	start := time.Now()

	collectFromIndex(&out, &root, os.Args[2], 0, byte(maxDistance))

	elapsed := time.Now().Sub(start)

	for key, value := range out {
		fmt.Println(key, value)
	}

	fmt.Println(elapsed)
}
