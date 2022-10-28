package processors

import (
	"bufio"
	"errors"
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

	var file *os.File = nil
	var err error = nil

	//TODO: валидация размера файла
	if _, err := os.Stat(processor.filename); errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(processor.filename)
	} else {
		file, err = os.Open(processor.filename)
	}

	if err != nil {
		return err
	}

	processor.Reader = bufio.NewReader(file)

	return nil
}

func (processor *FileProcessor) WriteFromChan(ch chan []byte) {

}
