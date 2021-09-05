package main

import (
	"os"

	"github.com/andlabs/ui"
)

const (
	applicationName = "dupeGuru Reference Batch Modifier"
)

var (
	result   = &Result{}
	fileInfo os.FileInfo
)

func main() {
	ui.Main(gui)
}
