package xls

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWBOpenFile(t *testing.T) {
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

func TestWBSummaryInfo(t *testing.T) {
	testCases := []struct {
		fName      string
		title      string
		subject    string
		author     string
		keywords   string
		comment    string
		lastAuthor string
		appName    string
		category   string
		manager    string
		company    string
	}{
		{smallFile, "", "", "cleit", "", "", "leitiennec", "", "", "", ""},
		{bigFile, "Some report", "Energy", "", "energy", "comment 1", "", "", "", "", ""},
	}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName)

		assert.NoError(t, err)
		assert.NotNil(t, wb)

		assert.Equal(t, tc.title, wb.Summary.Title)
		assert.Equal(t, tc.subject, wb.Summary.Subject)
		assert.Equal(t, tc.author, wb.Summary.Author)
		assert.Equal(t, tc.keywords, wb.Summary.Keywords)
		assert.Equal(t, tc.comment, wb.Summary.Comment)
		assert.Equal(t, tc.lastAuthor, wb.Summary.LastAuthor)
		wb.Close()
	}
}

func TestWBMeta(t *testing.T) {
	testCases := []struct {
		fName       string
		activeSheet uint
		charset     string
		sheetCount  uint
	}{
		{fName: smallFile, activeSheet: 0, charset: "UTF-8", sheetCount: 3},
		{fName: bigFile, activeSheet: 3, charset: "UTF-8", sheetCount: 13},
	}

	for _, tc := range testCases {
		wb, err := OpenFile(tc.fName)

		assert.NoError(t, err)
		assert.NotNil(t, wb)

		assert.Equal(t, tc.activeSheet, wb.ActiveSheetIdx)
		assert.Equal(t, tc.charset, wb.Charset)
		assert.Equal(t, tc.sheetCount, wb.SheetCount)

		wb.Close()
	}

}
