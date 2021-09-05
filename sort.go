package main

import "path/filepath"

type SortByPath []File

func (a SortByPath) Len() int           { return len(a) }
func (a SortByPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByPath) Less(i, j int) bool { return a[i].Path < a[j].Path }

type SortByName []File

func (a SortByName) Len() int           { return len(a) }
func (a SortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByName) Less(i, j int) bool { return filepath.Base(a[i].Path) < filepath.Base(a[j].Path) }

type SortByExt []File

func (a SortByExt) Len() int           { return len(a) }
func (a SortByExt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByExt) Less(i, j int) bool { return filepath.Ext(a[i].Path) < filepath.Ext(a[j].Path) }
