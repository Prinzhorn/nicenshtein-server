package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

//A Trie structure that maps runes to a list of following (child-) runes.
type node struct {
	children map[rune]*node
	word     string
}

var root node = node{make(map[rune]*node), ""}

func indexFile(fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		nextWord := strings.TrimSpace(scanner.Text())
		indexWord(nextWord)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func indexWord(word string) {
	if len(word) == 0 {
		return
	}

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

func findWord(word string) bool {
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

func collectClosestWords(out *map[string]byte, currentNode *node, word string, distance byte, maxDistance byte) {
	//We have eated all runes, let's see if we have reached a node with a valid word.
	if len(word) == 0 {
		if currentNode.word != "" {
			knownDistance, ok := (*out)[currentNode.word]

			//We have not seen this word or we have found a smaller distance.
			if !ok || distance < knownDistance {
				(*out)[currentNode.word] = distance
			}
		}

		return
	}

	if distance < maxDistance {
		for runeValue, _ := range currentNode.children {
			//Substitution (replace the first rune with the current one).
			collectClosestWords(out, currentNode, string(runeValue)+word[1:], distance+1, maxDistance)

			//Insertion (add the current rune as prefix).
			collectClosestWords(out, currentNode, string(runeValue)+word, distance+1, maxDistance)
		}

		//Deletion (skip first rune).
		collectClosestWords(out, currentNode, word[1:], distance+1, maxDistance)
	}

	runeValue, _ := utf8.DecodeRuneInString(word)
	nextNode := currentNode.children[runeValue]

	if nextNode != nil {
		//Move forward by one rune without incrementing the distance.
		collectClosestWords(out, nextNode, word[1:], distance, maxDistance)
	}
}

func millionHandler(w http.ResponseWriter, req *http.Request) {
	word := req.URL.Path[len("/1e6"):]

	if word == "" {
		fmt.Fprintln(w, "Specify a word like /1e6/password")
		return
	}

	out := make(map[string]byte)

	start := time.Now()

	collectClosestWords(&out, &root, word, 0, 2)

	log.Printf("Searched %s in %s\n", word, time.Now().Sub(start))

	/*for word, distance := range out {
		fmt.Fprintln(w, word, distance)
	}*/

	jsonData, err := json.MarshalIndent(out, "", "  ")

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not convert to json", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	fmt.Fprintln(w, "Heelo")
}

func main() {
	start := time.Now()

	err := indexFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Indexed in %s\n", time.Now().Sub(start))

	http.HandleFunc("/1e6/", millionHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}
