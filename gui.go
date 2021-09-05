package main

import (
	"encoding/xml"
	"io/ioutil"
	"math"
	"os"

	"github.com/andlabs/ui"
)

func gui() {
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
