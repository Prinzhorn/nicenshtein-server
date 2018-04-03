package main

import (
	"encoding/json"
	"fmt"
	"github.com/prinzhorn/nicenshtein"
	"log"
	"net/http"
	"os"
	"time"
)

var nice = nicenshtein.NewNicenshtein()

func millionHandler(w http.ResponseWriter, req *http.Request) {
	word := req.URL.Path[len("/1e6/"):]

	if word == "" {
		http.Error(w, "Specify a word like /1e6/password", 400)
		return
	}

	start := time.Now()

	out := make(map[string]byte)
	nice.CollectClosestWords(&out, word, 2)

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

func defaultHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}

	fmt.Fprintln(w, "Heelo")
}

func main() {
	start := time.Now()

	err := nice.IndexFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
		return
	}

	log.Printf("Indexed in %s\n", time.Now().Sub(start))

	http.HandleFunc("/1e6/", millionHandler)
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8080", nil)
}
