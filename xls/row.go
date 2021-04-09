package xls

/*
#include <xls.h>

typedef struct st_row_data st_row_data;
typedef struct st_cell_data st_cell_data;
*/
import "C"
import (
	"fmt"
	"strconv"
	"unsafe"
)

type Row struct {
	scr   *C.st_row_data
	Index uint
	Cells []*Cell
}

func makeRow(wb *WorkBook, src *C.st_row_data, collCnt int) *Row {
	cellDataSz := unsafe.Sizeof(C.st_cell_data{})
	cells := make([]*Cell, collCnt)

	cellPtr := unsafe.Pointer(src.cells.cell)
	for i := 0; i < collCnt; i++ {
		cells[i] = makeCell(wb, (*C.st_cell_data)(cellPtr))
		cellPtr = unsafe.Pointer(uintptr(cellPtr) + cellDataSz)
	}

	return &Row{
		scr:   src,
		Index: uint(uint32(src.index)),
		Cells: cells,
	}
}

type recordType int

const (
	recordBlank    = 0x0201
	recordLabelSST = 0x00FD
	recordLabel    = 0x0204
	recordFormula  = 0x0006
	recordNumber   = 0x0203
	recordRK       = 0x027E
	recordMulRK    = 0x00BD
)

func makeCell(wb *WorkBook, src *C.st_cell_data) *Cell {
	id := recordType(src.id)
	str := C.GoString(src.str)

	cell := &Cell{
		Style:   &CellStyle{uint16(src.xf), wb.xf[uint16(src.xf)]},
		Hidden:  uint8(src.isHidden) != 0,
		Rowspan: uint16(src.rowspan),
		Colspan: uint16(src.colspan),
	}

	switch id {
	case recordBlank:
		cell.Value = &BlankValue{}
	case recordFormula:
		l := int64(src.l)
		if l == 0 {
			cell.Value = &FloatValue{
				Val: float64(src.d),
			}
		} else {
			if str == "bool" {
				var b bool
				// precision ?
				if float64(src.d) > 0 {
					b = true
				}

				cell.Value = &BoolValue{
					Val: b,
				}
			}
			if str == "error" {
				cell.Value = &ErrValue{
					Code: int32(src.d),
				}

			} else {
				cell.Value = &StringValue{
					Val: str,
				}
			}
		}
	case recordLabelSST, recordLabel:
		cell.Value = &StringValue{
			Val: str,
		}
	case recordNumber, recordRK, recordMulRK:
		cell.Value = &FloatValue{
			Val: float64(src.d),
		}
	default:
		cell.Value = &UnknownValue{
			Val: str,
		}
	}

	return cell
}

type CellValue interface {
	fmt.Stringer
}

type Cell struct {
	Style   *CellStyle
	Value   CellValue
	Hidden  bool
	Rowspan uint16
	Colspan uint16
}

type BlankValue struct {
}

func (b *BlankValue) String() string {
	return ""
}

type FloatValue struct {
	Val float64
}

func (f *FloatValue) String() string {
	return strconv.FormatFloat(f.Val, 'f', -1, 64)
}

type BoolValue struct {
	Val bool
}

func (b *BoolValue) String() string {
	return strconv.FormatBool(b.Val)
}

type ErrValue struct {
	Code int32
}

func (e *ErrValue) String() string {
	return strconv.FormatInt(int64(e.Code), 10)
}

type StringValue struct {
	Val string
}

func (s *StringValue) String() string {
	return s.Val
}

type UnknownValue struct {
	Val string
}

func (u *UnknownValue) String() string {
	return u.Val
}
