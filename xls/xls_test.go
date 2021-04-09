package xls

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const smallFile = "./../test/data/test_small.xls"
const bigFile = "./../test/data/test_big.xls"
const styleFile = "./../test/data/test_style.xls"

func TestLibVersion(t *testing.T) {
	version := LibVersion()

	assert.NotEqual(t, "", version)
}
