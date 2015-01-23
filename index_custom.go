package main

import (
	"bufio"
	"fmt"
	"github.com/blevesearch/bleve"
	"log"
	"os"
	"strings"
	"time"
)

type Event struct {
	UID         string    `json:"uid"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Speaker     string    `json:"speaker"`
	Location    string    `json:"location"`
	Category    string    `json:"category"`
	URL         string    `json:"url"`
	Start       time.Time `json:"start"`
	Duration    float64   `json:"duration"`
}

func main() {
	// START OMIT
	enFieldMapping := bleve.NewTextFieldMapping()
	enFieldMapping.Analyzer = "en"

	eventMapping := bleve.NewDocumentMapping()
	eventMapping.AddFieldMappingsAt("summary", enFieldMapping)
	eventMapping.AddFieldMappingsAt("description", enFieldMapping)

	kwFieldMapping := bleve.NewTextFieldMapping()
	kwFieldMapping.Analyzer = "keyword"

	eventMapping.AddFieldMappingsAt("url", kwFieldMapping)
	eventMapping.AddFieldMappingsAt("category", kwFieldMapping)

	mapping := bleve.NewIndexMapping()
	mapping.DefaultMapping = eventMapping

	index, err := bleve.New("custom.bleve", mapping)
	if err != nil {
		log.Fatal(err)
	}
	// END OMIT

	count := 0
	batch := bleve.NewBatch()
	for event := range parseEvents() {
		batch.Index(event.UID, event)
		if batch.Size() > 100 {
			err := index.Batch(batch)
			if err != nil {
				log.Fatal(err)
			}
			count += batch.Size()
			batch = bleve.NewBatch()
		}
	}
	if batch.Size() > 0 {
		index.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
		count += batch.Size()
	}
	fmt.Printf("Indexed %d Events\n", count)
}

const iCalTimeFormat = "20060102T150405"

func parseEvents() chan *Event {
	rv := make(chan *Event)

	go func() {
		defer close(rv)

		var event *Event
		file, err := os.Open("fosdem.ical.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "BEGIN") {
				event = new(Event)
			} else if strings.HasPrefix(line, "END") {
				rv <- event
			} else if strings.HasPrefix(line, "SUMMARY") {
				colon := strings.Index(line, ":")
				if colon > 0 {
					event.Summary = line[colon+1:]
				}
			} else if strings.HasPrefix(line, "DESCRIPTION") {
				colon := strings.Index(line, ":")
				if colon > 0 {
					desc := line[colon+1:]
					desc = strings.TrimSpace(desc)

					if strings.HasPrefix(desc, "<p>") {
						desc = desc[3:]
					}
					if strings.HasSuffix(desc, "</p>") {
						desc = desc[:len(desc)-4]
					}
					if len(desc) > 0 {
						event.Description = desc
					}
				}
			} else if strings.HasPrefix(line, "LOCATION") {
				colon := strings.Index(line, ":")
				if colon > 0 {
					location := line[colon+1:]
					location = strings.TrimSpace(location)
					event.Location = location
				}
			} else if strings.HasPrefix(line, "STATUS") {
				// ignore all CONFIRMED
			} else if strings.HasPrefix(line, "CLASS") {
				// ignore all PUBLIC
			} else if strings.HasPrefix(line, "TZID") {
				// ignore all Europe-Brussels
			} else if strings.HasPrefix(line, "CATEGORIES") {
				colon := strings.Index(line, ":")
				if colon > 0 {
					cat := line[colon+1:]
					cat = strings.TrimSpace(cat)
					event.Category = cat
				}
			} else if strings.HasPrefix(line, "URL") {
				colon := strings.Index(line, ":")
				if colon > 0 {
					url := line[colon+1:]
					url = strings.TrimSpace(url)
					event.URL = url
				}
			} else if strings.HasPrefix(line, "METHOD") {
				// ignore all PUBLISH
			} else if strings.HasPrefix(line, "UID") {
				colon := strings.Index(line, ":")
				if colon > 0 {
					uid := line[colon+1:]
					uid = strings.TrimSpace(uid)
					event.UID = uid
				}
			} else if strings.HasPrefix(line, "DTSTART") {
				colon := strings.Index(line, ":")
				if colon > 0 {
					start := line[colon+1:]
					start = strings.TrimSpace(start)
					startTime, err := time.Parse(iCalTimeFormat, start)
					if err == nil {
						event.Start = startTime
					}
				}
			} else if strings.HasPrefix(line, "DTEND") {
				colon := strings.Index(line, ":")
				if colon > 0 {
					end := line[colon+1:]
					end = strings.TrimSpace(end)
					endTime, err := time.Parse(iCalTimeFormat, end)
					if err == nil {
						if !event.Start.IsZero() {
							duration := endTime.Sub(event.Start)
							event.Duration = duration.Minutes()
						}
					}
				}
			} else if strings.HasPrefix(line, "ATTENDEE") {
				attendeeParts := strings.Split(line, ";")
				for _, part := range attendeeParts {
					if strings.HasPrefix(part, "CN") {
						equal := strings.Index(part, "=")
						if equal > 0 {
							cn := part[equal+1:]
							cn = strings.TrimSpace(cn)
							if strings.HasSuffix(cn, "\":invalid:nomail") {
								cn = cn[:len(cn)-len("\":invalid:nomail")]
							}
							if strings.HasPrefix(cn, "\"") {
								cn = cn[1:]
							}
							event.Speaker = cn
						}
					}
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}()

	return rv
}

func init() {
	os.RemoveAll("custom.bleve")
}
