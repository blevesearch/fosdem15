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

	tq1 := bleve.NewTermQuery("text")                       // HLQUERY
	tq2 := bleve.NewTermQuery("search")                     // HLQUERY
	q := bleve.NewConjunctionQuery([]bleve.Query{tq1, tq2}) // HLQUERY
	req := bleve.NewSearchRequest(q)
	req.Highlight = bleve.NewHighlightWithStyle("html")
	req.Fields = []string{"summary", "speaker"}
	res, err := index.Search(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

// END OMIT
