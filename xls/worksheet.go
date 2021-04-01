package xls

/*
#include <stdio.h>
#include <stdlib.h>
#include <xls.h>
*/
import "C"

type WorkSheet struct {
	src *C.xlsWorkSheet

	Name string
}

func (ws *WorkSheet) Close() error {
	C.xls_close_WS(ws.src)
	return nil
}
