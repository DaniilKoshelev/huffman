package tree

import (
	"errors"
	"huffman/src/core/bitsbuffer"
)

func Restore(buffer *bitsbuffer.Buffer, nodesCount uint8) (*Tree, error) {
	tree := newTree()
	tree.nodesCount = nodesCount

	err := tree.extractNodes(buffer)

	if err != nil {
		return nil, err
	}

	tree.buildTree() // можно строить код сразу при парсинге файла, не использовать тогда эту функцию

	return tree, nil
}

func (tree *Tree) extractNodes(buffer *bitsbuffer.Buffer) error {
	if buffer == nil {
		return errors.New("buffer is not set")
	}

	root, err := tree.extractNode(buffer)

	tree.nodes.PushBack(root)

	if err != nil {
		return err
	}

	return nil
}

func (tree *Tree) extractNode(buffer *bitsbuffer.Buffer) (node, error) {
	if tree.alreadyReadNodes == tree.nodesCount {
		return nil, nil
	}

	bit, err := buffer.ReadBit()

	if err != nil {
		return nil, errors.New("error: could not read next tree bit")
	}

	if bit == 1 {
		newByte, err := buffer.ReadByte()

		if err != nil {
			return nil, errors.New("error: could not read word for initial node")
		}

		tree.alreadyReadNodes++
		newWord := &word{newByte, 0, nil}
		tree.Words[newByte] = newWord

		newNode := newInitialNode()
		newNode.setWord(newWord)

		return newNode, nil
	}

	left, err := tree.extractNode(buffer)

	if err != nil {
		return nil, err
	}

	right, err := tree.extractNode(buffer)

	if err != nil {
		return nil, err
	}

	return newAbstractNode(left, right), nil
}
