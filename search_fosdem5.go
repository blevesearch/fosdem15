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

	longTalk := 110.0                               // HLQUERY
	q := bleve.NewNumericRangeQuery(&longTalk, nil) // HLQUERY
	req := bleve.NewSearchRequest(q)
	req.Highlight = bleve.NewHighlightWithStyle("html")
	req.Fields = []string{"summary", "speaker", "duration"}
	res, err := index.Search(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

// END OMIT
