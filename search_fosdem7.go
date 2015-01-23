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

	qString := `+description:text `       // HLQUERY
	qString += `summary:"text indexing" ` // HLQUERY
	qString += `summary:believe~2 `       // HLQUERY
	qString += `-description:lucene `     // HLQUERY
	qString += `duration:>30`             // HLQUERY
	q := bleve.NewQueryStringQuery(qString)
	req := bleve.NewSearchRequest(q)
	req.Highlight = bleve.NewHighlightWithStyle("html")
	req.Fields = []string{"summary", "speaker", "description", "duration"}
	res, err := index.Search(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

// END OMIT
