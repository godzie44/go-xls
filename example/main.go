package main

import (
	"bytes"
	"fmt"
	"github.com/godzie44/go-xls/xls"
	"log"
)

func main() {
	//open xls workbook from file
	wb, err := xls.OpenFile("./test/data/test_style.xls", "UTF-8")
	if err != nil {
		log.Fatal(err)
	}
	defer wb.Close()

	// open first worksheet
	sheet, err := wb.OpenWorkSheet(0)
	if err != nil {
		log.Fatal(err)
	}
	defer sheet.Close()

	html := bytes.Buffer{}

	// set styles from workbook
	html.WriteString(fmt.Sprintf("<style type=\"text/css\">\n%s</style>\n", wb.CSS()))

	html.WriteString("<table border=0 cellspacing=0 cellpadding=2>")
	for _, row := range sheet.Rows {
		html.WriteString("<tr>")
		for _, cell := range row.Cells {
			if cell.Hidden {
				continue
			}

			html.WriteString("<td")
			if cell.Colspan != 0 {
				html.WriteString(fmt.Sprintf(" colspan=%d", cell.Colspan))
			}
			if cell.Rowspan != 0 {
				html.WriteString(fmt.Sprintf(" rowspan=%d", cell.Rowspan))
			}

			html.WriteString(fmt.Sprintf(" class=%s", cell.Style.CSSClass()))

			html.WriteString(">")

			html.WriteString(cell.Value.String())

			html.WriteString("</td>")
		}

		html.WriteString("</tr>")
	}
	html.WriteString("</table>")

	fmt.Println(html.String())
}
