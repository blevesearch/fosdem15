package main

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
)

// START OMIT
func main() {

	index, err := bleve.Open("fosdem.bleve")
	if err != nil {
		log.Fatal(err)
	}

	q := bleve.NewTermQuery("bleve") // HL
	req := bleve.NewSearchRequest(q)
	req.Highlight = bleve.NewHighlightWithStyle("html") // HL
	req.Fields = []string{"summary", "speaker"}
	res, err := index.Search(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

// END OMIT
