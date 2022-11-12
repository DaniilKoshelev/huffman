package flow

import (
	"huffman/src/core"
	"huffman/src/processors"
)

type Flow struct {
	*processors.FlagProcessor
	inputFileProcessor  *processors.FileProcessor
	outputFileProcessor *processors.FileProcessor
}

func NewFlow() *Flow {
	return &Flow{}
}

func (flow *Flow) Init() error {
	flow.FlagProcessor = processors.NewFlagProcessor()

	err := flow.ProcessInput()

	flow.inputFileProcessor = processors.NewFileProcessor(flow.InputFilename)
	flow.outputFileProcessor = processors.NewFileProcessor(flow.OutputFilename)

	if err != nil {
		return err
	}

	if flow.IsEncodeMode() {
		return flow.startEncode()
	}

	return flow.startDecode()
}

func (flow *Flow) startEncode() error {
	err := flow.inputFileProcessor.OpenFileToRead()

	if err != nil {
		return err
	}

	defer flow.inputFileProcessor.CloseFile() // TODO: fixme

	encoder := core.NewEncoder()

	if err != nil {
		return err
	}

	err = flow.outputFileProcessor.OpenFileToWrite()

	defer flow.outputFileProcessor.CloseFile() // TODO: fixme

	if err != nil {
		return err
	}

	err = encoder.Init(flow.inputFileProcessor.Reader)

	if err != nil {
		return err
	}

	err = flow.inputFileProcessor.ResetCursor()

	if err != nil {
		return err
	}

	err = encoder.Encode(flow.inputFileProcessor.Reader, flow.outputFileProcessor.Writer)

	if err != nil {
		return err
	}

	return nil
}

func (flow *Flow) startDecode() error {
	err := flow.inputFileProcessor.OpenFileToRead()

	if err != nil {
		return err
	}

	defer flow.inputFileProcessor.CloseFile() // TODO: fixme

	decoder := core.NewDecoder()

	if err != nil {
		return err
	}

	err = flow.outputFileProcessor.OpenFileToWrite()

	defer flow.outputFileProcessor.CloseFile() // TODO: fixme

	if err != nil {
		return err
	}

	err = decoder.Init(flow.inputFileProcessor.Reader)

	if err != nil {
		return err
	}

	err = flow.inputFileProcessor.ResetCursor()

	if err != nil {
		return err
	}

	err = decoder.Decode(flow.inputFileProcessor.Reader, flow.outputFileProcessor.Writer)

	if err != nil {
		return err
	}

	return nil
}
