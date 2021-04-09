package xls

import (
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"testing"
)

type rowTS struct {
	suite.Suite
	openFn func(fName string, charset string) (*WorkBook, error)
}

func (suite *rowTS) TestRowsExtract() {
	testCases := []struct {
		fName       string
		SheetNumber int
		RowCount    int
	}{
		{smallFile, 0, 24},
		{bigFile, 0, 7321},
	}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName, "UTF-8")
		suite.NoError(err)

		ws, err := wb.OpenWorkSheet(tc.SheetNumber)
		suite.NoError(err)
		suite.NotNil(ws)

		suite.Len(ws.Rows, tc.RowCount)

		ws.Close()
		wb.Close()
	}
}

func (suite *rowTS) TestCellValues() {
	testCases := []struct {
		fName    string
		sheetNum int
		cases    []struct {
			Row    int
			Col    int
			Type   interface{}
			StrVal string
		}
	}{
		{fName: smallFile, sheetNum: 0, cases: []struct {
			Row    int
			Col    int
			Type   interface{}
			StrVal string
		}{
			{0, 0, &StringValue{}, "a"},
			{2, 1, &BlankValue{}, ""},
			{1, 0, &FloatValue{}, "2"},
			{16, 0, &FloatValue{}, "17"},
			{12, 1, &StringValue{}, "This is a horizontally merged cell"},
			{22, 1, &StringValue{}, "\"doubly quoted string\""},
			{22, 2, &StringValue{}, "'singly quoted string'"},
			{23, 2, &StringValue{}, "The quick brown fox"},
			{23, 3, &FloatValue{}, "2"},
		}},
		{fName: bigFile, sheetNum: 0, cases: []struct {
			Row    int
			Col    int
			Type   interface{}
			StrVal string
		}{
			{3, 1, &StringValue{}, "Асбест "},
			{6731, 1, &StringValue{}, "Алейская"},
			{3, 4, &FloatValue{}, "1031.123"},
			{6731, 4, &FloatValue{}, "1081.357"},
			{7320, 0, &ErrValue{}, "29"},
		}},
		{fName: bigFile, sheetNum: 11, cases: []struct {
			Row    int
			Col    int
			Type   interface{}
			StrVal string
		}{
			{3, 3, &StringValue{}, "Свердловская область"},
			{6731, 3, &StringValue{}, "Алтайский край"},
			{3, 2, &FloatValue{}, "110"},
			{6731, 2, &FloatValue{}, "110"},
		}},
	}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName, "UTF-8")
		suite.NoError(err)

		ws, err := wb.OpenWorkSheet(tc.sheetNum)
		suite.NoError(err)
		suite.NotNil(ws)

		for _, c := range tc.cases {
			cell := ws.Rows[c.Row].Cells[c.Col]

			suite.IsType(c.Type, cell.Value)

			switch v := cell.Value.(type) {
			case *FloatValue:
				suite.True(strings.Contains(v.String(), c.StrVal))
			default:
				suite.Equal(c.StrVal, v.String())
			}
		}

		ws.Close()
		wb.Close()
	}
}

func TestRowSuite(t *testing.T) {
	suite.Run(t, &rowTS{
		openFn: OpenFile,
	})
	suite.Run(t, &rowTS{
		openFn: func(fName string, charset string) (*WorkBook, error) {
			data, err := os.ReadFile(fName)
			if err != nil {
				return nil, err
			}
			return Open(data, charset)
		},
	})
}
