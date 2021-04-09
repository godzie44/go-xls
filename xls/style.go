package xls

/*
#include <stdlib.h>
#include <xls.h>

typedef struct st_font_data st_font_data;
typedef struct st_xf_data st_xf_data;
*/
import "C"
import (
	"fmt"
	"unsafe"
)

var defaultFont = Font{
	Height:     200,
	Flag:       0,
	color:      0,
	Bold:       0,
	Escapement: 0,
	Underline:  0,
	Family:     0,
	Charset:    0,
	Name:       "Arial",
}

type Font struct {
	Height     uint16
	Flag       uint16
	color      C.WORD
	Bold       uint16
	Escapement uint16
	Underline  byte
	Family     byte
	Charset    byte
	Name       string
}

func (f *Font) Color() uint32 {
	return uint32(C.xls_getColor(f.color, 0))
}

func parseFontTable(src *C.xlsWorkBook) []*Font {
	fontCnt := uint32(src.fonts.count)

	fonts := make([]*Font, 0, fontCnt)
	fontDataPtr := unsafe.Pointer(src.fonts.font)
	fontDataSz := unsafe.Sizeof(C.st_font_data{})
	for i := 0; i < int(fontCnt); i++ {
		f := (*C.st_font_data)(fontDataPtr)
		fonts = append(fonts, &Font{
			Height:     uint16(f.height),
			Flag:       uint16(f.flag),
			color:      f.color,
			Bold:       uint16(f.bold),
			Escapement: uint16(f.escapement),
			Underline:  byte(f.underline),
			Family:     byte(f.family),
			Charset:    byte(f.charset),
			Name:       C.GoString(f.name),
		})

		fontDataPtr = unsafe.Pointer(uintptr(fontDataPtr) + fontDataSz)
	}

	return fonts
}

type xf struct {
	Font        *Font
	Format      uint16
	Type        uint16
	Align       byte
	Rotation    byte
	Ident       byte
	UsedAttr    byte
	LineStyle   uint32
	LineColor   uint32
	groundColor uint16
}

// GroundColor background color (RGB).
func (x *xf) GroundColor() uint32 {
	return uint32(C.xls_getColor(C.WORD(x.groundColor&0x7f), 0))
}

func parseXFTable(src *C.xlsWorkBook, fontTbl []*Font) []*xf {
	xfCnt := uint32(src.xfs.count)

	xfs := make([]*xf, 0, xfCnt)
	xfDataPtr := unsafe.Pointer(src.xfs.xf)
	xfDataSz := unsafe.Sizeof(C.st_xf_data{})
	for i := 0; i < int(xfCnt); i++ {
		xfData := (*C.st_xf_data)(xfDataPtr)

		font, fontIDX := &defaultFont, uint16(xfData.font)
		if fontIDX > 0 {
			font = fontTbl[fontIDX-1]
		}

		xfs = append(xfs, &xf{
			Font:        font,
			Format:      uint16(xfData.format),
			Type:        uint16(xfData._type),
			Align:       byte(xfData.align),
			Rotation:    byte(xfData.rotation),
			Ident:       byte(xfData.ident),
			UsedAttr:    byte(xfData.usedattr),
			LineStyle:   uint32(xfData.linestyle),
			LineColor:   uint32(xfData.linecolor),
			groundColor: uint16(xfData.groundcolor),
		})

		xfDataPtr = unsafe.Pointer(uintptr(xfDataPtr) + xfDataSz)
	}

	return xfs
}

type CellStyle struct {
	id uint16
	*xf
}

// CSSClass return cell css class name.
func (c *CellStyle) CSSClass() string {
	return fmt.Sprintf("xf%d", c.id)
}

// CSS return style table as list of css classes.
func (wb *WorkBook) CSS() string {
	css := C.xls_getCSS(wb.src)
	defer C.free(unsafe.Pointer(css))

	return C.GoString(css)
}
