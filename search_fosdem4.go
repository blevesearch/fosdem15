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

	tq1 := bleve.NewTermQuery("text")
	tq2 := bleve.NewTermQuery("search")
	tq3 := bleve.NewTermQuery("believe") // HLBLEVE
	q := bleve.NewConjunctionQuery(
		[]bleve.Query{tq1, tq2, tq3})
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
