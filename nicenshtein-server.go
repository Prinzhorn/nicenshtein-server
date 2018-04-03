package main

import (
	"encoding/json"
	"github.com/prinzhorn/nicenshtein"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var niceWords, nicePasswords nicenshtein.Nicenshtein
var indexHTML []byte

func niceHandler(w http.ResponseWriter, req *http.Request, nice *nicenshtein.Nicenshtein, word string) {
	start := time.Now()

	out := make(map[string]byte)
	nice.CollectWords(&out, word, 2)

	log.Printf("Searched %s in %s\n", word, time.Now().Sub(start))

	jsonData, err := json.MarshalIndent(out, "", "  ")

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not convert to json", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func englishHandler(w http.ResponseWriter, req *http.Request) {
	word := req.URL.Path[len("/english/"):]

	if word == "" {
		http.Error(w, "Specify a word as path", 400)
		return
	}

	niceHandler(w, req, &niceWords, word)
}

func passwordsHandler(w http.ResponseWriter, req *http.Request) {
	word := req.URL.Path[len("/passwords/"):]

	if word == "" {
		http.Error(w, "Specify a word as path", 400)
		return
	}

	niceHandler(w, req, &nicePasswords, word)
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	w.Write(indexHTML)
}

func main() {
	var err error

	indexHTML, err = ioutil.ReadFile("./index.html")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	start := time.Now()

	niceWords = nicenshtein.NewNicenshtein()
	err = niceWords.IndexFile("./data/words.txt")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	log.Printf("Indexed English words in %s\n", time.Now().Sub(start))

	start = time.Now()

	nicePasswords = nicenshtein.NewNicenshtein()
	err = nicePasswords.IndexFile("./data/10-million-password-list-top-1000000.txt")

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	log.Printf("Indexed passwords in %s\n", time.Now().Sub(start))

	http.HandleFunc("/english/", englishHandler)
	http.HandleFunc("/passwords/", passwordsHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}
