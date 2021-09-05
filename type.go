package main

import (
	"encoding/xml"
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

func (r *Result) locate(row, column int) (int, int, File) {
	for i, v := range r.Groups {
		if count := len(v.Files); row < count {
			return i, row, v.Files[row]
		} else {
			row -= len(v.Files)
		}
	}
	return 0, 0, File{}
}
