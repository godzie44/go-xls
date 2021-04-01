package xls

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenFile(t *testing.T) {
	testCases := []struct {
		fName string
	}{{smallFile}, {bigFile}}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName)

		assert.NoError(t, err)
		assert.NotNil(t, wb)

		wb.Close()
	}
}

func TestSheetCount(t *testing.T) {
	testCases := []struct {
		fName string
		cnt   uint
	}{{smallFile, 3}, {bigFile, 24}}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName)
		assert.NoError(t, err)
		assert.NotNil(t, wb)

		assert.Equal(t, tc.cnt, wb.SheetCount)

		wb.Close()
	}
}
