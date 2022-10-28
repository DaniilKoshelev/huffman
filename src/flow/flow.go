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
	flow.inputFileProcessor = processors.NewFileProcessor()
	flow.outputFileProcessor = processors.NewFileProcessor()

	err := flow.ProcessInput()

	if err != nil {
		return err
	}

	if flow.IsEncodeMode() {
		return flow.startEncode()
	}

	return flow.startDecode()
}

func (flow *Flow) startEncode() error {
	err := flow.inputFileProcessor.OpenFile(flow.InputFilename)

	if err != nil {
		return err
	}

	encoder := core.NewEncoder(flow.inputFileProcessor.Reader)
	_ = encoder.BuildTree()

	err = flow.outputFileProcessor.OpenFile(flow.OutputFilename)

	if err != nil {
		return err
	}

	outChan := make(chan []byte) // TODO: check buffered or not

	go encoder.Encode(outChan)
	go flow.outputFileProcessor.WriteFromChan(outChan)

	return nil
}

func (flow *Flow) startDecode() error {
	// TODO:
	return nil
}
