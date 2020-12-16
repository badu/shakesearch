package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	searcher := Searcher{}

	if err := searcher.Load("completeworks.txt"); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks    string
	CompleteWorksLen int
	SuffixArray      *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(strings.ToLower(query[0]))
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)

		if err := enc.Encode(results); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error loading filename %q: %w", filename, err)
	}
	s.CompleteWorks = string(dat)
	dat = bytes.ToLower(dat)
	s.SuffixArray = suffixarray.New(dat)
	s.CompleteWorksLen = len(dat)
	return nil
}

func (s *Searcher) Search(query string) []string {
	var results []string
	now := time.Now()
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	for _, idx := range idxs {
		fragment := ""
		if idx < 250 {
			fragment = strings.ReplaceAll(s.CompleteWorks[:idx+250], query, "<b>"+query+"</b>")

		} else if idx+250 > len(s.CompleteWorks) {
			fragment = strings.ReplaceAll(s.CompleteWorks[idx-250:], query, "<b>"+query+"</b>")
		} else {
			fragment = strings.ReplaceAll(s.CompleteWorks[idx-250:idx+250], query, "<b>"+query+"</b>")
		}
		results = append(results, fragment)
	}
	fmt.Printf("search took %d ns.\n", time.Now().Sub(now).Nanoseconds())
	if len(results) == 0 {
		fmt.Printf("no results found.\n")
	}
	return results
}
