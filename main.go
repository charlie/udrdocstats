package main

import (
	"os"
	"fmt"
	"encoding/xml"
	"log"
	"io"
	"strings"
)

var stopwords = []string{"wsrcd", "subinfo", "wsinfo", "chrginfo", "timeinfo"}

func include(elementName string) bool {
	for _, stopword := range stopwords {
		if elementName == stopword {
			return false
		}
	}
	return true
}

func printElementNames(file *os.File) {
	seqmap := make(map[string]int)
	d := xml.NewDecoder(file)
	buffer := make([]string, 0)
	d.Token() // discard <?xml ...>
	d.Token() // discard <DespatchFile>
	i := 0
	for {
		token, err := d.Token()
		if token == nil && err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error processing %s: %v\n", file.Name(), err)
		}
		if start, ok := token.(xml.StartElement); ok {
			if include(start.Name.Local) {
				buffer = append(buffer, start.Name.Local)
			}
		} else if end, ok := token.(xml.EndElement); ok {
			if end.Name.Local == "wsrcd" {
				key := strings.Join(buffer, "-")
				count := seqmap[key]
				seqmap[key] = count + 1
				buffer = buffer[:0]
				i += 1
				if i % 5000 == 0 {
					//fmt.Printf("%d records processed\n", i)
				}
			}
		}
	}
	fmt.Printf("Processed %d records\n", i)
	for sequence, count := range seqmap {
		fmt.Printf("%05d %s\n", count, sequence)
	}

}

func main() {
	for _, file := range os.Args[1:] {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		printElementNames(f)
	}
}
