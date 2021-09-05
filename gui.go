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
	sortByFileExtButtonText          = "Sort by File Extension Name"
	sortBySizeButtonText             = "Sort by Size"
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

	defaultAction := ui.NewHorizontalBox()
	defaultAction.SetPadded(true)
	defaultAction.Disable()
	outer.Append(defaultAction, false)

	fetchFileInfo := ui.NewButton(fetchFileInfoButtonText)
	defaultAction.Append(fetchFileInfo, true)

	sortByFilePath := ui.NewButton(sortByFilePathButtonText)
	defaultAction.Append(sortByFilePath, true)

	sortByFileName := ui.NewButton(sortByFileNameButtonText)
	defaultAction.Append(sortByFileName, true)

	sortByFileExt := ui.NewButton(sortByFileExtButtonText)
	defaultAction.Append(sortByFileExt, true)

	fileInfoAction := ui.NewHorizontalBox()
	fileInfoAction.SetPadded(true)
	fileInfoAction.Disable()
	outer.Append(fileInfoAction, false)

	sortBySize := ui.NewButton(sortBySizeButtonText)
	fileInfoAction.Append(sortBySize, true)

	sortByCreationTime := ui.NewButton(sortByCreationTimeButtonText)
	fileInfoAction.Append(sortByCreationTime, true)

	sortByModificationTime := ui.NewButton(sortByModificationTimeButtonText)
	fileInfoAction.Append(sortByModificationTime, true)

	sortByAccessName := ui.NewButton(sortByAccessNameButtonText)
	fileInfoAction.Append(sortByAccessName, true)

	tableModel := ui.NewTableModel(result)
	table := ui.NewTable(&ui.TableParams{
		Model: tableModel,
	})
	table.AppendTextColumn("Group", 0, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Path", 1, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Size", 2, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Creation Time", 3, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Modification Time", 4, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Access Time", 5, ui.TableModelColumnNeverEditable, nil)
	table.AppendTextColumn("Words", 6, ui.TableModelColumnNeverEditable, nil)
	table.AppendCheckboxColumn("Reference", 7, ui.TableModelColumnNeverEditable)
	table.AppendCheckboxColumn("Marked", 8, ui.TableModelColumnNeverEditable)
	table.AppendButtonColumn("Set Reference", 9, ui.TableModelColumnAlwaysEditable)
	table.AppendButtonColumn("Mark", 10, ui.TableModelColumnAlwaysEditable)
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
			defaultAction.Enable()
			fileInfoAction.Disable()
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
		result.fetchFileInfo()
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
		fileInfoAction.Enable()
	})

	sortByFilePath.OnClicked(func(*ui.Button) {
		if strings.Contains(sortByFilePath.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortByPath(w.Files)))
			}
			sortByFilePath.SetText(sortByFilePathButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(SortByPath(w.Files))
			}
			sortByFilePath.SetText(sortByFilePathButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	sortByFileName.OnClicked(func(*ui.Button) {
		if strings.Contains(sortByFileName.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortByName(w.Files)))
			}
			sortByFileName.SetText(sortByFileNameButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(SortByName(w.Files))
			}
			sortByFileName.SetText(sortByFileNameButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	sortByFileExt.OnClicked(func(*ui.Button) {
		if strings.Contains(sortByFileExt.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortByExt(w.Files)))
			}
			sortByFileExt.SetText(sortByFileExtButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(SortByExt(w.Files))
			}
			sortByFileExt.SetText(sortByFileExtButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	sortBySize.OnClicked(func(*ui.Button) {
		if strings.Contains(sortBySize.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortBySize(w.Files)))
			}
			sortBySize.SetText(sortBySizeButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(SortBySize(w.Files))
			}
			sortBySize.SetText(sortBySizeButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	sortByCreationTime.OnClicked(func(*ui.Button) {
		if strings.Contains(sortByCreationTime.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortByCreation(w.Files)))
			}
			sortByCreationTime.SetText(sortByCreationTimeButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(SortByCreation(w.Files))
			}
			sortByCreationTime.SetText(sortByCreationTimeButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	sortByModificationTime.OnClicked(func(*ui.Button) {
		if strings.Contains(sortByModificationTime.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortByModification(w.Files)))
			}
			sortByModificationTime.SetText(sortByModificationTimeButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(SortByModification(w.Files))
			}
			sortByModificationTime.SetText(sortByModificationTimeButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	sortByAccessName.OnClicked(func(*ui.Button) {
		if strings.Contains(sortByAccessName.Text(), ascending) {
			for _, w := range result.Groups {
				sort.Sort(sort.Reverse(SortByAccess(w.Files)))
			}
			sortByAccessName.SetText(sortByAccessNameButtonText + " " + descending)
		} else {
			for _, w := range result.Groups {
				sort.Sort(SortByAccess(w.Files))
			}
			sortByAccessName.SetText(sortByAccessNameButtonText + " " + ascending)
		}
		rows := result.NumRows(tableModel)
		updateTable(rows, rows)
	})

	window.Show()
}
