package xls

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const smallFile = "./../test/data/test_small.xls"
const bigFile = "./../test/data/test_big.xls"

func TestLibVersion(t *testing.T) {
	version := LibVersion()

	assert.Equal(t, "1.6.2", version)
}