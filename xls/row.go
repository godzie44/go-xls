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
	Cells []Cell
}

func makeRow(wb *WorkBook, src *C.st_row_data, collCnt int) *Row {
	cellDataSz := unsafe.Sizeof(C.st_cell_data{})
	cells := make([]Cell, collCnt)

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

func makeCell(wb *WorkBook, src *C.st_cell_data) Cell {
	id := recordType(src.id)
	str := C.GoString(src.str)

	bc := baseCell{uint16(src.xf), wb.xf[uint16(src.xf)]}

	switch id {
	case recordBlank:
		return &BlankCell{
			baseCell: bc,
		}
	case recordFormula:
		l := int64(src.l)
		if l == 0 {
			return &FloatCell{
				baseCell: bc,
				Val:      float64(src.d),
			}
		} else {
			if str == "bool" {
				var b bool
				// precision ?
				if float64(src.d) > 0 {
					b = true
				}
				return &BoolCell{
					baseCell: bc,
					Val:      b,
				}
			}
			if str == "error" {
				return &ErrCell{
					baseCell: bc,
					Code:     int32(src.d),
				}
			}
			return &StringCell{
				baseCell: bc,
				Val:      str,
			}
		}
	case recordLabelSST, recordLabel:
		return &StringCell{
			Val:      str,
			baseCell: bc,
		}
	case recordNumber, recordRK, recordMulRK:
		return &FloatCell{
			baseCell: bc,
			Val:      float64(src.d),
		}
	default:
		return &UnknownCell{
			baseCell: bc,
			Val:      str,
		}
	}
}

type baseCell struct {
	styleID uint16
	style   *xf
}

func (b *baseCell) Style() *CellStyle {
	return &CellStyle{b.styleID, b.style}
}

type Cell interface {
	fmt.Stringer
	Style() *CellStyle
}

type BlankCell struct {
	baseCell
}

func (b *BlankCell) String() string {
	return ""
}

type FloatCell struct {
	baseCell
	Val float64
}

func (f *FloatCell) String() string {
	return strconv.FormatFloat(f.Val, 'f', -1, 64)
}

type BoolCell struct {
	baseCell
	Val bool
}

func (b *BoolCell) String() string {
	return strconv.FormatBool(b.Val)
}

type ErrCell struct {
	baseCell
	Code int32
}

func (e *ErrCell) String() string {
	return strconv.FormatInt(int64(e.Code), 10)
}

type StringCell struct {
	baseCell
	Val string
}

func (s *StringCell) String() string {
	return s.Val
}

type UnknownCell struct {
	baseCell
	Val string
}

func (u *UnknownCell) String() string {
	return u.Val
}
