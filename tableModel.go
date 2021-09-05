package main

import (
	"strconv"
	"time"

	"github.com/andlabs/ui"
)

func (r *Result) ColumnTypes(tm *ui.TableModel) []ui.TableValue {
	return []ui.TableValue{
		ui.TableString(""), // column  0 Group Index
		ui.TableString(""), // column  1 Path
		ui.TableString(""), // column  2 Size
		ui.TableString(""), // column  3 Creation Time
		ui.TableString(""), // column  4 Modification Time
		ui.TableString(""), // column  5 Access Time
		ui.TableString(""), // column  6 Words
		ui.TableInt(0),     // column  7 IsRef
		ui.TableInt(0),     // column  8 Marked
		ui.TableString(""), // column  9 IsRef button text
		ui.TableString(""), // column 10 Marked button text
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
		if !file.fetched {
			return ui.TableString("")
		} else {
			return ui.TableString(strconv.Itoa(int(file.size)))
		}
	case 3:
		if !file.fetched {
			return ui.TableString("")
		} else {
			return ui.TableString(file.time.creation.Format(time.RFC3339Nano))
		}
	case 4:
		if !file.fetched {
			return ui.TableString("")
		} else {
			return ui.TableString(file.time.modification.Format(time.RFC3339Nano))
		}
	case 5:
		if !file.fetched {
			return ui.TableString("")
		} else {
			return ui.TableString(file.time.access.Format(time.RFC3339Nano))
		}
	case 6:
		return ui.TableString(file.Words)
	case 7:
		if file.IsRef == "y" {
			return ui.TableTrue
		} else {
			return ui.TableFalse
		}
	case 8:
		if file.Marked == "y" {
			return ui.TableTrue
		} else {
			return ui.TableFalse
		}
	case 9:
		if file.IsRef == "y" {
			return ui.TableString("Unset reference")
		} else {
			return ui.TableString("Set as reference")
		}
	case 10:
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
	case 8:
		switch file.IsRef {
		case "n":
			r.Groups[groupId].Files[fileId].IsRef = "y"
		case "y":
			r.Groups[groupId].Files[fileId].IsRef = "n"
		}
	case 9:
		switch file.Marked {
		case "n":
			r.Groups[groupId].Files[fileId].Marked = "y"
		case "y":
			r.Groups[groupId].Files[fileId].Marked = "n"
		}
	}
	tm.RowChanged(row)
}
