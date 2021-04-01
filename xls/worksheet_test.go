package xls

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenWorkSheet(t *testing.T) {
	testCases := []struct {
		fName             string
		existsSheetNum    int
		notExistsSheetNum int
	}{
		{smallFile, 0, 3},
		{bigFile, 3, 25},
	}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName)
		assert.NoError(t, err)

		ws, err := wb.OpenWorkSheet(tc.existsSheetNum)
		assert.NoError(t, err)
		assert.NotNil(t, ws)
		ws.Close()

		ws, err = wb.OpenWorkSheet(tc.notExistsSheetNum)
		assert.ErrorIs(t, err, errInvalidWorkSheetNumber)
		assert.Nil(t, ws)

		wb.Close()
	}
}

func TestWorkSheetName(t *testing.T) {
	testCases := []struct {
		fName                           string
		firstSheetName, secondSheetName string
	}{
		{smallFile, "Sheet1", "Sheet2"},
		{bigFile, "0", "1"},
	}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName)
		assert.NoError(t, err)

		ws, err := wb.OpenWorkSheet(0)
		assert.NoError(t, err)
		assert.NotNil(t, ws)
		assert.Equal(t, tc.firstSheetName, ws.Name)
		ws.Close()

		ws, err = wb.OpenWorkSheet(1)
		assert.NoError(t, err)
		assert.NotNil(t, ws)
		assert.Equal(t, tc.secondSheetName, ws.Name)
		ws.Close()

		wb.Close()
	}
}
