package xls

/*
#include <xls.h>

typedef struct st_row_data st_row_data;
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// OpenWorkSheet opens sheet by its number, starts from 0.
func (wb *WorkBook) OpenWorkSheet(number int) (*WorkSheet, error) {
	cWS := C.xls_getWorkSheet(wb.src, C.int(number))
	if cWS == nil {
		return nil, errInvalidWorkSheetNumber
	}

	cErr := C.xls_parseWorkSheet(cWS)
	err := libXLSErr(cErr).intoErr()
	if err != nil {
		return nil, fmt.Errorf("work sheet %d: %w", number, err)
	}

	return &WorkSheet{src: cWS, Name: wb.sheetNames[number], Rows: parseRows(wb, cWS)}, nil
}

func parseRows(wb *WorkBook, cWS *C.xlsWorkSheet) []*Row {
	rowCnt := int(uint16(cWS.rows.lastrow)) + 1
	rowDataSz := unsafe.Sizeof(C.st_row_data{})
	rows := make([]*Row, rowCnt)

	rowPtr := unsafe.Pointer(cWS.rows.row)
	for i := 0; i < rowCnt; i++ {
		rows[i] = makeRow(wb, (*C.st_row_data)(rowPtr), int(uint16(cWS.rows.lastcol))+1)
		rowPtr = unsafe.Pointer(uintptr(rowPtr) + rowDataSz)
	}

	return rows
}

type WorkSheet struct {
	src *C.xlsWorkSheet

	Name string

	Rows []*Row
}

func (ws *WorkSheet) Close() {
	C.xls_close_WS(ws.src)
}
