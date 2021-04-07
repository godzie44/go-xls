package xls

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestRowsExtract(t *testing.T) {
	testCases := []struct {
		fName       string
		SheetNumber int
		RowCount    int
	}{
		{smallFile, 0, 24},
		{bigFile, 0, 7321},
	}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName)
		assert.NoError(t, err)

		ws, err := wb.OpenWorkSheet(tc.SheetNumber)
		assert.NoError(t, err)
		assert.NotNil(t, ws)

		assert.Len(t, ws.Rows, tc.RowCount)

		ws.Close()
		wb.Close()
	}
}

func TestCellValues(t *testing.T) {
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
			{0, 0, &StringCell{}, "a"},
			{2, 1, &BlankCell{}, ""},
			{1, 0, &FloatCell{}, "2"},
			{16, 0, &FloatCell{}, "17"},
			{12, 1, &StringCell{}, "This is a horizontally merged cell"},
			{22, 1, &StringCell{}, "\"doubly quoted string\""},
			{22, 2, &StringCell{}, "'singly quoted string'"},
			{23, 2, &StringCell{}, "The quick brown fox"},
			{23, 3, &FloatCell{}, "2"},
		}},
		{fName: bigFile, sheetNum: 0, cases: []struct {
			Row    int
			Col    int
			Type   interface{}
			StrVal string
		}{
			{3, 1, &StringCell{}, "Асбест "},
			{6731, 1, &StringCell{}, "Алейская"},
			{3, 4, &FloatCell{}, "1031.123"},
			{6731, 4, &FloatCell{}, "1081.357"},
			{7320, 0, &ErrCell{}, "29"},
		}},
		{fName: bigFile, sheetNum: 11, cases: []struct {
			Row    int
			Col    int
			Type   interface{}
			StrVal string
		}{
			{3, 3, &StringCell{}, "Свердловская область"},
			{6731, 3, &StringCell{}, "Алтайский край"},
			{3, 2, &FloatCell{}, "110"},
			{6731, 2, &FloatCell{}, "110"},
		}},
	}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName)
		assert.NoError(t, err)

		ws, err := wb.OpenWorkSheet(tc.sheetNum)
		assert.NoError(t, err)
		assert.NotNil(t, ws)

		for _, c := range tc.cases {
			cell := ws.Rows[c.Row].Cells[c.Col]

			assert.IsType(t, c.Type, cell)

			if _, isFloat := cell.(*FloatCell); isFloat {
				assert.True(t, strings.Contains(cell.String(), c.StrVal))
			} else {
				assert.Equal(t, c.StrVal, cell.String())
			}
		}

		ws.Close()
		wb.Close()
	}
}
