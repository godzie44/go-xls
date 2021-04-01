package xls

/*
#include <stdio.h>
#include <stdlib.h>
#include <xls.h>
*/
import "C"

func LibVersion() string {
	v := C.xls_getVersion()
	return C.GoString(v)
}
