package utils

import (
	"fmt"
	"os"
	"testing"
)

func TestVErrorAt(t *testing.T) {
	fmt.Fprintln(os.Stderr, "TestVErrorAt")
}
