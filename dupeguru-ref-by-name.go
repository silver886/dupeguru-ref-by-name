package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type File struct {
	Path   string `xml:"path,attr"`
	Words  string `xml:"words,attr"`
	IsRef  string `xml:"is_ref,attr"`
	Marked string `xml:"marked,attr"`
}
type Match struct {
	First      string `xml:"first,attr"`
	Second     string `xml:"second,attr"`
	Percentage string `xml:"percentage,attr"`
}
type Group struct {
	Files  []File  `xml:"file"`
	Matchs []Match `xml:"match"`
}
type Result struct {
	XMLName xml.Name `xml:"results"`
	Groups  []Group  `xml:"group"`
}

func main() {
	for i, v := range os.Args {
		if i == 0 {
			continue
		}
		fmt.Println("Creating backup . . .")

		src, err := os.Open(v)
		if err != nil {
			log.Fatalln("Cannot open file.")
		}

		dest, err := os.Create(v + ".bak")
		if err != nil {
			log.Fatalln("Cannot create backup.")
		}

		_, err = io.Copy(dest, src)
		if err != nil {
			log.Fatalln("Cannot backup file.")
		}

		src.Close()
		dest.Close()

		before, err := ioutil.ReadFile(v)
		if err != nil {
			log.Fatalln("Cannot read file content.")
		}

		result := Result{}
		if xml.Unmarshal(before, &result) != nil {
			log.Fatalln("Cannot parse file content.")
		}

		for _, w := range result.Groups {
			files := []string{}
			for _, x := range w.Files {
				files = append(files, x.Path)
			}
			fmt.Println(files)
			sort.Strings(files)
			fmt.Println(files)
			for i, x := range files {
				w.Files[i].Path = x
			}
		}

		after, err := xml.Marshal(result)
		if err != nil {
			log.Fatalln("Cannot encode sorted content.")
		}

		if os.WriteFile(v, after, os.ModeAppend) != err {
			log.Fatalln("Cannot write sorted content.")
		}
	}
}
