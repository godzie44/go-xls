package xls

/*
#include <stdio.h>
#include <stdlib.h>
#include <xls.h>


typedef struct st_sheet_data st_sheet_data;
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type libXLSErr int

const (
	noError libXLSErr = iota
	errorOpen
	errorSeek
	errorRead
	errorParse
	errorMalloc
)

var (
	errOpenFile               = errors.New("open file error")
	errReadFile               = errors.New("open file error")
	errParseFile              = errors.New("parsing file error")
	errInvalidWorkSheetNumber = errors.New("invalid worksheet number")
)

func (l libXLSErr) IntoErr() error {
	switch l {
	case noError:
		return nil
	case errorOpen:
		return errOpenFile
	case errorSeek:
		return errors.New("unknown error")
	case errorRead:
		return errReadFile
	case errorParse:
		return errParseFile
	case errorMalloc:
		return errors.New("unknown error")
	default:
		return errors.New("unknown error")
	}
}

type WorkBook struct {
	src *C.xlsWorkBook

	SheetCount uint
	sheetNames []string
}

func parseWorkBook(src *C.xlsWorkBook) (*WorkBook, error) {
	cErr := C.xls_parseWorkBook(src)
	err := libXLSErr(cErr).IntoErr()
	if err != nil {
		return nil, fmt.Errorf("work book: %w", err)
	}

	sheetCount := uint(src.sheets.count)

	firstSheetDataPtr := unsafe.Pointer(src.sheets.sheet)
	sheetDataSz := unsafe.Sizeof(C.st_sheet_data{})
	sheetNames := make([]string, sheetCount)
	for i := 0; i < int(sheetCount); i++ {
		sheetDataPtr := unsafe.Pointer(uintptr(firstSheetDataPtr) + (uintptr(i) * sheetDataSz))

		sheetData := *(*C.st_sheet_data)(sheetDataPtr)
		sheetNames[i] = C.GoString(sheetData.name)
	}

	return &WorkBook{
		src:        src,
		SheetCount: sheetCount,
		sheetNames: sheetNames,
	}, nil
}

func OpenFile(name string) (*WorkBook, error) {
	encoding := C.CString("UTF-8")
	defer C.free(unsafe.Pointer(encoding))

	f := C.CString(name)
	defer C.free(unsafe.Pointer(f))

	cErr := new(C.xls_error_t)
	wb := C.xls_open_file(f, encoding, cErr)

	err := libXLSErr(*cErr).IntoErr()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}

	return parseWorkBook(wb)
}

func (wb *WorkBook) Close() error {
	C.xls_close_WB(wb.src)
	return nil
}
