package tree

import (
	"bufio"
	"errors"
	"io"
	"sort"
)

func Create(reader *bufio.Reader) (*Tree, error) {
	tree := newTree()
	err := tree.countWords(reader)

	if err != nil {
		return nil, err
	}

	tree.pushInitialNodes()
	tree.compressTree()
	tree.buildTree()

	return tree, nil
}

func (tree *Tree) countWords(reader *bufio.Reader) error {
	if reader == nil {
		return errors.New("reader is not set")
	}

	for {
		newByte, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		curWord := tree.Words[newByte]

		if curWord != nil {
			curWord.count++
		} else {
			tree.Words[newByte] = &word{newByte, 1, nil}
		}
	}

	return nil
}

func (tree *Tree) pushInitialNodes() {
	var words []*word

	for _, word := range tree.Words {
		if word != nil {
			words = append(words, word)
		}
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].count < words[j].count
	})

	for _, word := range words {
		newNode := newInitialNodeCount(word.count)
		newNode.setWord(word)

		tree.nodes.PushBack(newNode)
	}
}

func (tree *Tree) compressTree() {
	for tree.nodes.Len() != 1 {
		leftElement := tree.nodes.Front()
		rightElement := leftElement.Next()

		left := leftElement.Value.(node)
		right := rightElement.Value.(node)

		newNode := newAbstractNodeCount(left.getCount() + right.getCount())
		newNode.left = left
		newNode.right = right

		tree.nodes.Remove(leftElement)
		tree.nodes.Remove(rightElement)

		for e := tree.nodes.Front(); e != nil; e = e.Next() {
			element := e.Value.(node)

			if element.getCount() >= newNode.getCount() {
				tree.nodes.InsertBefore(newNode, e)

				break
			}

			if e.Next() == nil {
				tree.nodes.PushBack(newNode)

				break
			}
		}

		if tree.nodes.Len() == 0 {
			tree.nodes.PushBack(newNode)
		}
	}
}
