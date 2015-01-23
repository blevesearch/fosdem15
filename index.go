package main

import (
	"fmt"
	"log"
	"os"
)

// START OMIT
import "github.com/blevesearch/bleve" // HLIMPORT

type Person struct { // HLMODEL
	Name string // HLMODEL
} // HLMODEL

func main() {
	mapping := bleve.NewIndexMapping()               // HLMAPPING
	index, err := bleve.New("people.bleve", mapping) // HLNEW
	if err != nil {
		log.Fatal(err)
	}

	person := Person{"Marty Schoch"} // HLINDEX
	err = index.Index("m1", person)  // HLINDEX
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Indexed Document")
}

// END OMIT

func init() {
	os.RemoveAll("people.bleve")
}
