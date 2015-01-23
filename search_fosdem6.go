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

	lateSunday := "2015-02-01T17:30:00Z"           // HLQUERY
	q := bleve.NewDateRangeQuery(&lateSunday, nil) // HLQUERY
	q.SetField("start")                            // HLQUERY
	req := bleve.NewSearchRequest(q)
	req.Highlight = bleve.NewHighlightWithStyle("html")
	req.Fields = []string{"summary", "speaker", "start"}
	res, err := index.Search(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

// END OMIT
