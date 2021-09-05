package main

import (
	"strconv"

	"github.com/andlabs/ui"
)

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
