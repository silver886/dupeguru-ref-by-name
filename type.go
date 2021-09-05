package main

import (
	"encoding/xml"
	"os"
	"sync"
	"syscall"
	"time"
)

type File struct {
	Path    string `xml:"path,attr"`
	Words   string `xml:"words,attr"`
	IsRef   string `xml:"is_ref,attr"`
	Marked  string `xml:"marked,attr"`
	fetched bool
	size    int64
	time    struct {
		creation     time.Time
		modification time.Time
		access       time.Time
	}
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

func (r *Result) setFirstReference() {
	for i, v := range result.Groups {
		for j := range v.Files {
			if j == 0 {
				r.Groups[i].Files[j].IsRef = "y"
			} else {
				r.Groups[i].Files[j].IsRef = "n"
			}
		}
	}
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

func (r *Result) fetchFileInfo() {
	var wg sync.WaitGroup
	for i, v := range r.Groups {
		wg.Add(len(v.Files))
		for j, w := range v.Files {
			groupId, fileId, file := i, j, w
			go func() {
				defer wg.Done()
				if fileInfo, err := os.Stat(file.Path); err != nil {
					return
				} else {
					r.Groups[groupId].Files[fileId].fetched = true
					r.Groups[groupId].Files[fileId].size = fileInfo.Size()
					fileStat := fileInfo.Sys().(*syscall.Stat_t)
					r.Groups[groupId].Files[fileId].time.creation = timespecToTime(fileStat.Ctimespec)
					r.Groups[groupId].Files[fileId].time.modification = timespecToTime(fileStat.Mtimespec)
					r.Groups[groupId].Files[fileId].time.access = timespecToTime(fileStat.Atimespec)
				}
			}()
		}
	}
	wg.Wait()
}
