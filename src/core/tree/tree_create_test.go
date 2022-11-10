package tree

import (
	"bufio"
	"bytes"
	"testing"
)

type buildTreeTest struct {
	bytes string
	words map[byte]uint16
}

var buildTreeTests = []buildTreeTest{
	{"a", map[byte]uint16{'a': 0}},
	{"ab", map[byte]uint16{'a': 32768, 'b': 0}},
	{"abc", map[byte]uint16{'a': 16384, 'b': 0, 'c': 32768}},
	{"FGGCCCZZZZ", map[byte]uint16{'F': 24576, 'G': 16384, 'C': 0, 'Z': 32768}},
}

func TestTreeCreate(t *testing.T) {
	for _, test := range buildTreeTests {
		reader := bytes.NewReader([]byte(test.bytes))
		bufferedReader := bufio.NewReader(reader)
		createdTree, _ := Create(bufferedReader)

		for encoderByte, encoderWord := range createdTree.words {
			if encoderWord == nil {
				continue
			}
			actualCode := encoderWord.code.Bits()
			expectedCode := test.words[byte(encoderByte)]

			if expectedCode != actualCode {
				t.Errorf("expected %q, actual %q", expectedCode, actualCode)
			}
		}
	}
}
