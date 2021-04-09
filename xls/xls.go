// Package bufio implements structs and functions for working with xls files.
package xls

/*
#include <stdio.h>
#include <stdlib.h>
#include <xls.h>
*/
import "C"

// LibVersion return current libxls version.
func LibVersion() string {
	v := C.xls_getVersion()
	return C.GoString(v)
}
