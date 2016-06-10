package parser

import (
	"fmt"
	"testing"
)

func TestReadNTParser(t *testing.T) {
	parser := NTParser{}
	for result := range parser.Read("test.nt") {
		fmt.Println(result)
	}
}
