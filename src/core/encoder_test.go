package core

import (
	"bufio"
	"bytes"
	"testing"
)

type buildTreeTest struct {
	bytes string
	words map[byte]string
}

var buildTreeTests = []buildTreeTest{
	{"a", map[byte]string{'a': "0"}},
	{"ab", map[byte]string{'a': "1", 'b': "0"}},
	{"abc", map[byte]string{'a': "01", 'b': "00", 'c': "1"}},
	{"FGGCCCZZZZ", map[byte]string{'F': "011", 'G': "010", 'C': "00", 'Z': "1"}},
}

func TestEncoderBuildTree(t *testing.T) {
	for _, test := range buildTreeTests {
		reader := bytes.NewReader([]byte(test.bytes))
		bufferedReader := bufio.NewReader(reader)
		encoder := NewEncoder(bufferedReader)

		_ = encoder.BuildTree()

		for encoderByte, encoderWord := range encoder.words {
			if encoderWord == nil {
				continue
			}
			actualCode := encoderWord.code.String()
			expectedCode := test.words[byte(encoderByte)]

			if expectedCode != actualCode {
				t.Errorf("expected %q, actual %q", expectedCode, actualCode)
			}
		}
	}
}

//type EncodeTest struct {
//	bytes        string
//	encodedBytes string
//}
//
//var EncodeTests = []EncodeTest{
//	{"FGGCCCZZZZ", ""}, 	// 01101001 00000001 111
//}
//
//func TestEncoderEncode(t *testing.T) {
//
//}
