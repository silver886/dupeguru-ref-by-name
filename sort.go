package main

import (
	"path/filepath"
)

type SortByPath []File

func (s SortByPath) Len() int      { return len(s) }
func (s SortByPath) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortByPath) Less(i, j int) bool {
	return s[i].Path < s[j].Path
}

type SortByName []File

func (s SortByName) Len() int      { return len(s) }
func (s SortByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortByName) Less(i, j int) bool {
	return filepath.Base(s[i].Path) < filepath.Base(s[j].Path)
}

type SortByExt []File

func (s SortByExt) Len() int      { return len(s) }
func (s SortByExt) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortByExt) Less(i, j int) bool {
	return filepath.Ext(s[i].Path) < filepath.Ext(s[j].Path)
}

type SortBySize []File

func (s SortBySize) Len() int      { return len(s) }
func (s SortBySize) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortBySize) Less(i, j int) bool {
	return s[i].size < s[j].size
}

type SortByCreation []File

func (s SortByCreation) Len() int      { return len(s) }
func (s SortByCreation) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortByCreation) Less(i, j int) bool {
	return s[i].time.creation.Before(s[j].time.creation)
}

type SortByModification []File

func (s SortByModification) Len() int      { return len(s) }
func (s SortByModification) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortByModification) Less(i, j int) bool {
	return s[i].time.modification.Before(s[j].time.modification)
}

type SortByAccess []File

func (s SortByAccess) Len() int      { return len(s) }
func (s SortByAccess) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortByAccess) Less(i, j int) bool {
	return s[i].time.access.Before(s[j].time.access)
}
