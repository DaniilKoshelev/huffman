package processors

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type FileProcessor struct {
	Reader   *bufio.Reader
	filename string
}

func NewFileProcessor(filename string) *FileProcessor {
	processor := new(FileProcessor)

	processor.filename = filename

	return processor
}

func (processor *FileProcessor) OpenFile() error {
	if processor.filename == "" {
		return errors.New("filename is not set")
	}

	//TODO: валидация размера файла + проверка что файл существует
	file, err := os.Open(processor.filename)
	processor.Reader = bufio.NewReader(file)

	defer file.Close() // TODO: обработать закрытие

	if err != nil {
		return errors.New(fmt.Sprintf("error: %s\n", err.Error()))
	}

	return nil
}

func (processor *FileProcessor) WriteFromChan(ch chan []byte) {

}
