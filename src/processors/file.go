package processors

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type FileProcessor struct {
	Reader *bufio.Reader
}

func NewFileProcessor() *FileProcessor {
	return &FileProcessor{}
}

func (processor *FileProcessor) OpenFile(filename string) error {
	//TODO: валидация размера файла
	file, err := os.Open(filename)
	processor.Reader = bufio.NewReader(file)

	defer file.Close() // TODO: обработать закрытие

	if err != nil {
		return errors.New(fmt.Sprintf("error: %s\n", err.Error()))
	}

	return nil
}

func (processor *FileProcessor) WriteFromChan(ch chan []byte) {

}
