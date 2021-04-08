package xls

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStyleFont(t *testing.T) {
	testCases := []struct {
		row, col  int
		height    uint16
		name      string
		bold      uint16
		underline byte
		color     uint32
	}{
		{1, 0, 200, "Arial", 400, 0, 0xFF0000},
		{1, 1, 200, "Arial", 400, 1, 0x000000},
		{1, 2, 200, "Arial", 400, 0, 0xFFFFFF},
		{1, 3, 720, "Arial", 700, 0, 0x333399},
		{2, 0, 240, "FreeSerif", 400, 0, 0x000000},
		{2, 1, 320, "Symbola", 400, 0, 0x993366},
		{2, 2, 200, "Arial", 700, 0, 0x000000},
		{2, 3, 720, "Arial", 400, 0, 0x333399},
	}

	wb, err := OpenFile(styleFile)
	assert.NoError(t, err)
	defer wb.Close()

	sheet, err := wb.OpenWorkSheet(0)
	assert.NoError(t, err)
	defer sheet.Close()

	for _, tc := range testCases {
		st := sheet.Rows[tc.row].Cells[tc.col].Style()

		assert.Equal(t, tc.height, st.Font.Height)
		assert.Equal(t, tc.name, st.Font.Name)
		assert.Equal(t, tc.bold, st.Font.Bold)
		assert.Equal(t, tc.underline, st.Font.Underline)
		assert.Equal(t, tc.color, st.Font.Color())
	}
}

func TestStyle(t *testing.T) {
	testCases := []struct {
		row, col    int
		align       byte
		rotation    byte
		groundColor uint32
		cssClass    string
	}{
		{1, 0, 0x20, 0, 0, "xf21"},
		{1, 1, 0x20, 0, 0xCCCCFF, "xf22"},
		{1, 2, 0x22, 0, 0, "xf23"},
		{1, 3, 0x20, 0, 0, "xf24"},
		{2, 0, 0x20, 0, 0xFFCC00, "xf26"},
		{2, 1, 0x20, 0, 0x339966, "xf27"},
		{2, 2, 0x23, 0, 0, "xf28"},
		{2, 3, 0x20, 0x5a, 0, "xf29"},
	}

	wb, err := OpenFile(styleFile)
	assert.NoError(t, err)
	defer wb.Close()

	sheet, err := wb.OpenWorkSheet(0)
	assert.NoError(t, err)
	defer sheet.Close()

	for _, tc := range testCases {
		st := sheet.Rows[tc.row].Cells[tc.col].Style()

		assert.Equal(t, tc.align, st.Align)
		assert.Equal(t, tc.rotation, st.Rotation)
		assert.Equal(t, tc.groundColor, st.GroundColor())
		assert.Equal(t, tc.cssClass, st.CSSClass())
	}
}

func TestStyleTableToCSS(t *testing.T) {
	wb, err := OpenFile(styleFile)
	assert.NoError(t, err)
	defer wb.Close()

	var expectedCSS = `.xf0{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf1{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf2{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf3{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf4{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf5{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf6{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf7{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf8{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf9{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf10{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf11{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf12{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf13{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf14{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf15{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf16{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf17{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf18{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf19{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf20{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf21{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#FF0000;}
.xf22{ font-size:10pt;font-family: "Arial";background:#CCCCFF;text-align:left;vertical-align:bottom;color:#000000;text-decoration: underline;}
.xf23{ font-size:10pt;font-family: "Arial";background:#000000;text-align:center;vertical-align:bottom;color:#FFFFFF;}
.xf24{ font-size:36pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#333399;font-weight: bold;}
.xf25{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;}
.xf26{ font-size:12pt;font-family: "FreeSerif";background:#FFCC00;text-align:left;vertical-align:bottom;color:#000000;}
.xf27{ font-size:16pt;font-family: "Symbola";background:#339966;text-align:left;vertical-align:bottom;color:#993366;}
.xf28{ font-size:10pt;font-family: "Arial";background:#FFFFFF;text-align:right;vertical-align:bottom;color:#000000;font-weight: bold;}
.xf29{ font-size:36pt;font-family: "Arial";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#333399;}
.xf30{ font-size:12pt;font-family: "FreeSerif";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#000000;}
.xf31{ font-size:16pt;font-family: "Symbola";background:#FFFFFF;text-align:left;vertical-align:bottom;color:#993366;}
`

	assert.Equal(t, expectedCSS, wb.CSS())
}