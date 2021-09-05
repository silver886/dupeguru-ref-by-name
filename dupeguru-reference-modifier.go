package main

import (
	"os"

	"github.com/andlabs/ui"
)

const (
	linkApplicationName = "placeholder"
)

var (
	result   = &Result{}
	fileInfo os.FileInfo
)

func main() {
	ui.Main(gui)
}
