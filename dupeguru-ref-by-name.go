package main

import (
	"encoding/xml"
	"io/ioutil"
	"math"
	"os"
	"strconv"

	"github.com/andlabs/ui"
)

const (
	applicationName = "dupeGuru Reference Batch Modifier"
)

var (
	result   = &Result{}
	fileInfo os.FileInfo
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

func (r *Result) ColumnTypes(tm *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""), // column 0 Group Index
		ui.TableString(""), // column 1 Path
		ui.TableString(""), // column 2 Words
		ui.TableInt(0),     // column 3 IsRef
		ui.TableInt(0),     // column 4 Marked
		ui.TableString(""), // column 5 IsRef button text
		ui.TableString(""), // column 6 Marked button text
	}
}

func (r *Result) NumRows(tm *ui.TableModel) int {
	rowCount := 0
	for _, v := range r.Groups {
		rowCount += len(v.Files)
	}
	return rowCount
}

func (r *Result) CellValue(tm *ui.TableModel, row, column int) ui.TableValue {
	switch groupId, _, file := r.locate(row, column); column {
	case 0:
		return ui.TableString(strconv.Itoa(groupId + 1))
	case 1:
		return ui.TableString(file.Path)
	case 2:
		return ui.TableString(file.Words)
	case 3:
		if file.IsRef == "y" {
			return ui.TableTrue
		} else {
			return ui.TableFalse
		}
	case 4:
		if file.Marked == "y" {
			return ui.TableTrue
		} else {
			return ui.TableFalse
		}
	case 5:
		if file.IsRef == "y" {
			return ui.TableString("Unset reference")
		} else {
			return ui.TableString("Set as reference")
		}
	case 6:
		if file.Marked == "y" {
			return ui.TableString("Check")
		} else {
			return ui.TableString("Uncheck")
		}
	}
	panic("unreachable")
}

func (r *Result) SetCellValue(tm *ui.TableModel, row, column int, value ui.TableValue) {
	switch groupId, fileId, file := r.locate(row, column); column {
	case 5:
		switch file.IsRef {
		case "n":
			r.Groups[groupId].Files[fileId].IsRef = "y"
		case "y":
			r.Groups[groupId].Files[fileId].IsRef = "n"
		}
	case 6:
		switch file.Marked {
		case "n":
			r.Groups[groupId].Files[fileId].Marked = "y"
		case "y":
			r.Groups[groupId].Files[fileId].Marked = "n"
		}
	}
	tm.RowChanged(row)
}

func setupUI() {
	window := ui.NewWindow(applicationName, 1280, 720, true)
	window.SetMargined(true)

	outer := ui.NewVerticalBox()
	outer.SetPadded(true)
	window.SetChild(outer)

	file := ui.NewHorizontalBox()
	file.SetPadded(true)
	outer.Append(file, false)

	open := ui.NewButton("Open")
	file.Append(open, true)

	save := ui.NewButton("Save")
	save.Disable()
	file.Append(save, true)

	tableModel := ui.NewTableModel(result)
	table := ui.NewTable(&ui.TableParams{
		Model: tableModel,
	})
	table.AppendTextColumn("Group", 0, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Path", 1, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Words", 2, ui.TableModelColumnNeverEditable, nil)
	table.AppendCheckboxColumn("Reference", 3, ui.TableModelColumnNeverEditable)
	table.AppendCheckboxColumn("Marked", 4, ui.TableModelColumnNeverEditable)
	table.AppendButtonColumn("Set Reference", 5, ui.TableModelColumnAlwaysEditable)
	table.AppendButtonColumn("Mark", 6, ui.TableModelColumnAlwaysEditable)
	outer.Append(table, false)

	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	ui.OnShouldQuit(func() bool {
		window.Destroy()
		return true
	})

	open.OnClicked(func(*ui.Button) {
		currentRows := result.NumRows(tableModel)
		if filePath := ui.OpenFile(window); len(filePath) == 0 {
		} else if file, err := os.Open(filePath); err != nil {
			ui.MsgBoxError(window, "Open File Error", "Cannot open "+filePath)
		} else if info, err := file.Stat(); err != nil {
			ui.MsgBoxError(window, "Fetch File Info Error", "Cannot fetch info of "+filePath)
		} else if content, err := ioutil.ReadFile(filePath); err != nil {
			ui.MsgBoxError(window, "Read Content Error", "Cannot read "+filePath)
		} else if result.Groups = nil; xml.Unmarshal(content, result) != nil {
			ui.MsgBoxError(window, "Parse Content Error", "Cannot parse "+filePath)
		} else {
			fileInfo = info
			save.Enable()
			newRows := result.NumRows(tableModel)
			duplicated := int(math.Min(float64(currentRows), float64(newRows)))
			for i := 0; i < duplicated; i++ {
				tableModel.RowChanged(i)
			}
			if currentRows == duplicated {
				for i := duplicated; i < newRows; i++ {
					tableModel.RowInserted(i)
				}
			} else {
				for i := currentRows; i > duplicated; i-- {
					tableModel.RowDeleted(i - 1)
				}
			}
		}
	})

	save.OnClicked(func(*ui.Button) {
		if content, err := xml.Marshal(result); err != nil {
			ui.MsgBoxError(window, "Encode Content Error", "Cannot encode modified content")
		} else if filePath := ui.SaveFile(window); len(filePath) == 0 {
		} else if os.WriteFile(filePath, content, fileInfo.Mode()) != err {
			ui.MsgBoxError(window, "Write Content Error", "Cannot write to "+filePath)
		}
	})

	window.Show()
}

func main() {
	ui.Main(setupUI)
}
