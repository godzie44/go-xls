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
	errReadFile               = errors.New("read from file error")
	errParseFile              = errors.New("parsing file error")
	errSeekFile               = errors.New("seek within file error")
	errMalloc                 = errors.New("allocate memory error")
	errInvalidWorkSheetNumber = errors.New("invalid worksheet number")
	errUnknown                = errors.New("unknown error")
)

func (l libXLSErr) IntoErr() error {
	switch l {
	case noError:
		return nil
	case errorOpen:
		return errOpenFile
	case errorSeek:
		return errSeekFile
	case errorRead:
		return errReadFile
	case errorParse:
		return errParseFile
	case errorMalloc:
		return errMalloc
	default:
		return errUnknown
	}
}

type (
	WorkBook struct {
		src *C.xlsWorkBook

		Summary *SummaryInfo

		SheetCount     uint
		ActiveSheetIdx int
		Charset        string

		sheetNames []string
	}
	SummaryInfo struct {
		Title      string
		Subject    string
		Author     string
		Keywords   string
		Comment    string
		LastAuthor string
		AppName    string
		Category   string
		Manager    string
		Company    string
	}
)

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

	summary := C.xls_summaryInfo(src)
	defer C.xls_close_summaryInfo(summary)

	title := unsafe.Pointer(summary.title)
	subj := unsafe.Pointer(summary.subject)
	author := unsafe.Pointer(summary.author)
	keywords := unsafe.Pointer(summary.keywords)
	comment := unsafe.Pointer(summary.comment)
	lastAuthor := unsafe.Pointer(summary.lastAuthor)
	appName := unsafe.Pointer(summary.appName)
	category := unsafe.Pointer(summary.category)
	manager := unsafe.Pointer(summary.manager)
	company := unsafe.Pointer(summary.company)

	return &WorkBook{
		src:        src,
		SheetCount: sheetCount,
		sheetNames: sheetNames,
		Summary: &SummaryInfo{
			Title:      C.GoString((*C.char)(title)),
			Subject:    C.GoString((*C.char)(subj)),
			Author:     C.GoString((*C.char)(author)),
			Keywords:   C.GoString((*C.char)(keywords)),
			Comment:    C.GoString((*C.char)(comment)),
			LastAuthor: C.GoString((*C.char)(lastAuthor)),
			AppName:    C.GoString((*C.char)(appName)),
			Category:   C.GoString((*C.char)(category)),
			Manager:    C.GoString((*C.char)(manager)),
			Company:    C.GoString((*C.char)(company)),
		},
		ActiveSheetIdx: int(src.activeSheetIdx),
		Charset:        C.GoString(src.charset),
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
