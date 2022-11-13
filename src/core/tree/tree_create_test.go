package tree

import (
	"bufio"
	"bytes"
	"testing"
)

type buildTreeTest struct {
	bytes string
	words map[byte]uint64
}

var buildTreeTests = []buildTreeTest{
	{"a", map[byte]uint64{'a': 0}},
	{"ab", map[byte]uint64{'a': 1 << 63, 'b': 0}},
	{"abc", map[byte]uint64{'a': 1 << 62, 'b': 0, 'c': 1 << 63}},
	{"FGGCCCZZZZ", map[byte]uint64{'F': 1<<62 | 1<<61, 'G': 1 << 62, 'C': 0, 'Z': 1 << 63}},
}

func TestTreeCreate(t *testing.T) {
	for _, test := range buildTreeTests {
		reader := bytes.NewReader([]byte(test.bytes))
		bufferedReader := bufio.NewReader(reader)
		createdTree, _ := Create(bufferedReader)

		for encoderByte, encoderWord := range createdTree.Words {
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
