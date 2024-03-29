package processors

import (
	"bufio"
	"errors"
	"os"
)

type FileProcessor struct {
	Reader   *bufio.Reader
	Writer   *bufio.Writer
	file     *os.File
	filename string
}

func NewFileProcessor(filename string) *FileProcessor {
	processor := new(FileProcessor)

	processor.filename = filename

	return processor
}

func (processor *FileProcessor) OpenFileToRead() error {
	if processor.filename == "" {
		return errors.New("filename is not set")
	}

	//TODO: валидация размера файла
	file, err := os.Open(processor.filename)

	if err != nil {
		return err
	}

	processor.file = file
	processor.Reader = bufio.NewReader(file)

	return nil
}

// OpenFileToWrite TODO: объединить с методом выше
func (processor *FileProcessor) OpenFileToWrite() error {
	if processor.filename == "" {
		return errors.New("filename is not set")
	}

	//TODO: валидация размера файла
	file, err := os.OpenFile(processor.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		return err
	}

	processor.file = file
	processor.Writer = bufio.NewWriter(file)

	return nil
}

func (processor *FileProcessor) ResetCursor() error {
	_, err := processor.file.Seek(0, 0)

	if err != nil {
		return err
	}

	return nil
}

func (processor *FileProcessor) CloseFile() error {
	if processor.Writer != nil {
		err := processor.Writer.Flush()

		if err != nil {
			return err
		}
	}

	err := processor.file.Close()

	if err != nil {
		return err
	}

	return nil
}
