package main

import (
	"encoding/xml"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/andlabs/ui"
)

const (
	openButtonText                   = "Open"
	saveButtonText                   = "Save"
	fetchFileInfoButtonText          = "Fetch File Info"
	sortByFilePathButtonText         = "Sort by File Path"
	sortByFileNameButtonText         = "Sort by File Name"
	sortByCreationTimeButtonText     = "Sort by Creation Time"
	sortByModificationTimeButtonText = "Sort by Modification Time"
	sortByAccessNameButtonText       = "Sort by Access Name"

	ascending  = "▲"
	descending = "▼"
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

	open := ui.NewButton(openButtonText)
	file.Append(open, true)

	save := ui.NewButton(saveButtonText)
	save.Disable()
	file.Append(save, true)

	action := ui.NewHorizontalBox()
	action.SetPadded(true)
	outer.Append(action, false)

	fetchFileInfo := ui.NewButton(fetchFileInfoButtonText)
	fetchFileInfo.Disable()
	action.Append(fetchFileInfo, true)

	sortByFilePath := ui.NewButton(sortByFilePathButtonText)
	sortByFilePath.Disable()
	action.Append(sortByFilePath, true)

	sortByFileName := ui.NewButton(sortByFileNameButtonText)
	sortByFileName.Disable()
	action.Append(sortByFileName, true)

	sortByCreationTime := ui.NewButton(sortByCreationTimeButtonText)
	sortByCreationTime.Disable()
	action.Append(sortByCreationTime, true)

	sortByModificationTime := ui.NewButton(sortByModificationTimeButtonText)
	sortByModificationTime.Disable()
	action.Append(sortByModificationTime, true)

	sortByAccessName := ui.NewButton(sortByAccessNameButtonText)
	sortByAccessName.Disable()
	action.Append(sortByAccessName, true)

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
	updateTable := func(oldRows, newRows int) {
		duplicated := int(math.Min(float64(oldRows), float64(newRows)))
		for i := 0; i < duplicated; i++ {
			tableModel.RowChanged(i)
		}
		if oldRows == duplicated {
			for i := duplicated; i < newRows; i++ {
				tableModel.RowInserted(i)
			}
		} else {
			for i := oldRows; i > duplicated; i-- {
				tableModel.RowDeleted(i - 1)
			}
		}
	}

	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	ui.OnShouldQuit(func() bool {
		window.Destroy()
		return true
	})

	open.OnClicked(func(*ui.Button) {
		oldRows := result.NumRows(tableModel)
		if filePath := ui.OpenFile(window); len(filePath) == 0 {
			return
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
			fetchFileInfo.Enable()
			sortByFilePath.Enable()
			sortByFileName.Enable()
			updateTable(oldRows, result.NumRows(tableModel))
		}
	})

	save.OnClicked(func(*ui.Button) {
		if content, err := xml.Marshal(result); err != nil {
			ui.MsgBoxError(window, "Encode Content Error", "Cannot encode modified content")
		} else if filePath := ui.SaveFile(window); len(filePath) == 0 {
			return
		} else if os.WriteFile(filePath, content, fileInfo.Mode()) != err {
			ui.MsgBoxError(window, "Write Content Error", "Cannot write to "+filePath)
		}
	})

	fetchFileInfo.OnClicked(func(*ui.Button) {
	})

	sortByFilePath.OnClicked(func(*ui.Button) {
		if strings.Contains(sortByFilePath.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(SortByPath(w.Files))
			}
			sortByFilePath.SetText(sortByFilePathButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortByPath(w.Files)))
			}
			sortByFilePath.SetText(sortByFilePathButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	sortByFileName.OnClicked(func(*ui.Button) {
		if strings.Contains(sortByFileName.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(SortByName(w.Files))
			}
			sortByFileName.SetText(sortByFileNameButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortByName(w.Files)))
			}
			sortByFileName.SetText(sortByFileNameButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	sortByCreationTime.OnClicked(func(*ui.Button) {
	})

	sortByModificationTime.OnClicked(func(*ui.Button) {
	})

	sortByAccessName.OnClicked(func(*ui.Button) {
	})

	window.Show()
}
