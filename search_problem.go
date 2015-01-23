package main

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
)

func main() {

	index, err := bleve.Open("fosdem.bleve")
	if err != nil {
		log.Fatal(err)
	}

	// START OMIT
	q := bleve.NewTermQuery("haystack") // HL
	req := bleve.NewSearchRequest(q)
	req.Highlight = bleve.NewHighlightWithStyle("html")
	req.Fields = []string{"summary", "speaker"}
	res, err := index.Search(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	// END OMIT
}

func unused() *bleve.IndexMapping {
	// START2 OMIT
	mapping := bleve.NewIndexMapping() // HL
	// END2 OMIT
	return mapping
}
