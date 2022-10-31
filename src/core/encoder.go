package core

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
	"io"
	"sort"
)

type ensemble struct {
	words map[byte]*word
}

type word struct {
	count int64
	code  *bytes.Buffer
}

type Encoder struct {
	*ensemble
	tree   *list.List
	reader *bufio.Reader
}

func NewEncoder(reader *bufio.Reader) *Encoder {
	encoder := new(Encoder)
	encoder.reader = reader
	encoder.ensemble = new(ensemble)
	encoder.ensemble.words = make(map[byte]*word)
	encoder.tree = list.New()

	return encoder
}

func (encoder *Encoder) BuildTree() error {
	err := encoder.countWords()

	if err != nil {
		return err
	}

	encoder.pushInitialNodes()
	encoder.compressTree()
	encoder.walkTree()

	return nil
}

func (encoder *Encoder) countWords() error {
	if encoder.reader == nil {
		return errors.New("reader is not set")
	}

	for {
		newByte, err := encoder.reader.ReadByte()

		if err == io.EOF {
			break
		}

		curWord, exists := encoder.words[newByte]

		if exists {
			curWord.count++
		} else {
			encoder.words[newByte] = &word{1, nil}
		}
	}

	return nil
}

func (encoder *Encoder) Encode(ch chan []byte) {

}

func (encoder *Encoder) pushInitialNodes() {
	var words []*word

	for _, Word := range encoder.words {
		words = append(words, Word)
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].count < words[j].count
	})

	for _, word := range words {
		newNode := newInitialNode(word.count)
		newNode.setWord(word)

		encoder.tree.PushBack(newNode)
	}
}

func (encoder *Encoder) compressTree() {
	for encoder.tree.Len() != 1 { // TODO: check 0
		leftElement := encoder.tree.Front()
		rightElement := leftElement.Next()

		left := leftElement.Value.(node)
		right := rightElement.Value.(node)

		newNode := newAbstractNode(left.getCount() + right.getCount())
		newNode.Left = left
		newNode.Right = right

		encoder.tree.Remove(leftElement)
		encoder.tree.Remove(rightElement)

		for e := encoder.tree.Front(); e != nil; e = e.Next() {
			element := e.Value.(node)

			if element.getCount() >= newNode.getCount() {
				encoder.tree.InsertBefore(newNode, e)

				break
			}

			if e.Next() == nil {
				encoder.tree.PushBack(newNode)

				break
			}
		}

		if encoder.tree.Len() == 0 {
			encoder.tree.PushBack(newNode)
		}
	}
}

// TODO: рефакторинг нужен
func (encoder *Encoder) walkTree() {
	root := encoder.tree.Front().Value

	rootAbstract, ok := encoder.tree.Front().Value.(*abstractNode)
	code := new(bytes.Buffer)

	if !ok {
		code.WriteByte('0')
		root.(*initialNode).getWord().code = code

		return
	}

	left := rootAbstract.Left
	right := rootAbstract.Right

	if left != nil {
		leftCode := new(bytes.Buffer)
		leftCode.Write(code.Bytes())
		leftCode.WriteByte('1')
		encoder.walkFromNode(left, leftCode)
	}

	if right != nil {
		rightCode := new(bytes.Buffer)
		rightCode.Write(code.Bytes())
		rightCode.WriteByte('0')
		encoder.walkFromNode(right, rightCode)
	}
}

func (encoder *Encoder) walkFromNode(node node, code *bytes.Buffer) {
	_, ok := node.(*abstractNode)

	if !ok {
		node.(*initialNode).getWord().code = code

		return
	}

	left := node.(*abstractNode).Left
	right := node.(*abstractNode).Right

	if left != nil {
		leftCode := new(bytes.Buffer)
		leftCode.Write(code.Bytes())
		leftCode.WriteByte('1')
		encoder.walkFromNode(left, leftCode)
	}

	if right != nil {
		rightCode := new(bytes.Buffer)
		rightCode.Write(code.Bytes())
		rightCode.WriteByte('0')
		encoder.walkFromNode(right, rightCode)
	}
}

type node interface {
	getCount() int64
}

type abstractNode struct {
	count int64 // Количество раз, сколько встретился определенный байт
	Left  node
	Right node
}

type initialNode struct {
	count int64
	word  *word
}

func (node *initialNode) getWord() *word {
	return node.word
}

func (node *initialNode) setWord(word *word) {
	node.word = word
}

func newInitialNode(p int64) *initialNode {
	return &initialNode{count: p}
}

func newAbstractNode(p int64) *abstractNode {
	return &abstractNode{count: p}
}

func (node *initialNode) getCount() int64 {
	return node.count
}

func (node *abstractNode) getCount() int64 {
	return node.count
}
