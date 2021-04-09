package xls

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type worksheetTS struct {
	suite.Suite
	openFn func(fName string, charset string) (*WorkBook, error)
}

func (suite *worksheetTS) TestOpenWS() {
	testCases := []struct {
		fName             string
		existsSheetNum    int
		notExistsSheetNum int
	}{
		{smallFile, 0, 3},
		{bigFile, 3, 25},
	}

	for _, tc := range testCases {
		wb, err := suite.openFn(tc.fName, "UTF-8")
		suite.NoError(err)

		ws, err := wb.OpenWorkSheet(tc.existsSheetNum)
		suite.NoError(err)
		suite.NotNil(ws)
		ws.Close()

		ws, err = wb.OpenWorkSheet(tc.notExistsSheetNum)
		suite.ErrorIs(err, errInvalidWorkSheetNumber)
		suite.Nil(ws)

		wb.Close()
	}
}

func (suite *worksheetTS) TestWSName() {
	testCases := []struct {
		fName                           string
		firstSheetName, secondSheetName string
	}{
		{smallFile, "Sheet1", "Sheet2"},
		{bigFile, "0", "1"},
	}

	for _, tc := range testCases {
		wb, err := suite.openFn(tc.fName, "UTF-8")
		suite.NoError(err)

		ws, err := wb.OpenWorkSheet(0)
		suite.NoError(err)
		suite.NotNil(ws)
		suite.Equal(tc.firstSheetName, ws.Name)
		ws.Close()

		ws, err = wb.OpenWorkSheet(1)
		suite.NoError(err)
		suite.NotNil(ws)
		suite.Equal(tc.secondSheetName, ws.Name)
		ws.Close()

		wb.Close()
	}
}

func TestWorksheetSuite(t *testing.T) {
	suite.Run(t, &worksheetTS{
		openFn: OpenFile,
	})
	suite.Run(t, &worksheetTS{
		openFn: func(fName string, charset string) (*WorkBook, error) {
			data, err := os.ReadFile(fName)
			if err != nil {
				return nil, err
			}
			return Open(data, charset)
		},
	})
}
