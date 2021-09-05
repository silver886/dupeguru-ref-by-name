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
		ui.TableColor{},    // row background color
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
	switch groupId, group, fileId, file := r.locate(row, column); column {
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
			return ui.TableString("Uncheck")
		} else {
			return ui.TableString("Check")
		}
	case 11:
		if file.IsRef == "y" {
			return ui.TableColor{
				R: 0.5,
				G: 0.5,
				B: 0.5,
				A: 0.5,
			}
		} else if fileId == 0 {
			reference := true
			for i, v := range group.Files {
				if i == 0 {
					continue
				}
				if v.IsRef == "y" {
					reference = false
				}
			}
			if reference {
				return ui.TableColor{
					R: 0.5,
					G: 0.5,
					B: 0.5,
					A: 0.5,
				}
			}
		}
		return ui.TableColor{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
	}
	panic("unreachable")
}

func (r *Result) SetCellValue(tm *ui.TableModel, row, column int, value ui.TableValue) {
	switch groupId, group, fileId, file := r.locate(row, column); column {
	case 9:
		switch file.IsRef {
		case "n":
			r.Groups[groupId].Files[fileId].IsRef = "y"
		case "y":
			r.Groups[groupId].Files[fileId].IsRef = "n"
		}
		for i := row - fileId; i < row-fileId+len(group.Files); i++ {
			tm.RowChanged(i)
		}
	case 10:
		switch file.Marked {
		case "n":
			r.Groups[groupId].Files[fileId].Marked = "y"
		case "y":
			r.Groups[groupId].Files[fileId].Marked = "n"
		}
	}
	tm.RowChanged(row)
}
