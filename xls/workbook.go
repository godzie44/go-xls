package xls

/*
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
		ActiveSheetIdx uint
		Charset        string

		sheetNames []string

		xf []*xf
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

	sheetCnt := uint(uint32(src.sheets.count))

	sheetDataPtr := unsafe.Pointer(src.sheets.sheet)
	sheetDataSz := unsafe.Sizeof(C.st_sheet_data{})
	sheetNames := make([]string, sheetCnt)
	for i := 0; i < int(sheetCnt); i++ {
		sheetData := *(*C.st_sheet_data)(sheetDataPtr)
		sheetNames[i] = C.GoString(sheetData.name)

		sheetDataPtr = unsafe.Pointer(uintptr(sheetDataPtr) + sheetDataSz)
	}

	return &WorkBook{
		src:            src,
		SheetCount:     sheetCnt,
		sheetNames:     sheetNames,
		xf:             parseXFTable(src, parseFontTable(src)),
		Summary:        parseSummary(src),
		ActiveSheetIdx: uint(uint32(src.activeSheetIdx)),
		Charset:        C.GoString(src.charset),
	}, nil
}

func OpenFile(name string, charset string) (*WorkBook, error) {
	encoding := C.CString(charset)
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

func Open(buff []byte, charset string) (*WorkBook, error) {
	encoding := C.CString(charset)
	defer C.free(unsafe.Pointer(encoding))

	cErr := new(C.xls_error_t)
	wb := C.xls_open_buffer((*C.uchar)(&buff[0]), C.size_t(len(buff)), encoding, cErr)
	err := libXLSErr(*cErr).IntoErr()
	if err != nil {
		return nil, err
	}

	return parseWorkBook(wb)
}

func (wb *WorkBook) Close() {
	C.xls_close_WB(wb.src)
}

func parseSummary(src *C.xlsWorkBook) *SummaryInfo {
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

	return &SummaryInfo{
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
	}
}
