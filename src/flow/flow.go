package flow

import (
	"fmt"
	"huffman/src/core"
	"huffman/src/core/bitsbuffer"
	"huffman/src/entropy"
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
	} else if flow.IsEntropyMode() {
		return flow.startEntropy()
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

	fmt.Printf("Avg bits per word: %.2f\n", encoder.GetAverageBitsPerWord())

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

	inputFileBuffer := bitsbuffer.NewEmptyBuffer().SetIoReader(flow.inputFileProcessor.Reader)
	err = decoder.Init(inputFileBuffer)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	err = decoder.Decode(inputFileBuffer, flow.outputFileProcessor.Writer)

	if err != nil {
		return err
	}

	return nil
}

func (flow *Flow) startEntropy() error {
	err := flow.inputFileProcessor.OpenFileToRead()

	if err != nil {
		return err
	}

	defer flow.inputFileProcessor.CloseFile() // TODO: fixme

	entropyService := entropy.NewEntropy()
	entropyService.Init(flow.inputFileProcessor.Reader)

	E := entropyService.CalculateEntropy()
	E_XY := entropyService.CalculateConditionalEntropy()
	E_XYY := entropyService.CalculateDoubleConditionalEntropy()

	fmt.Printf("H(X): %.2f\n", E)
	fmt.Printf("H(X|X): %.2f\n", E_XY)
	fmt.Printf("H(X|XX): %.2f\n", E_XYY)

	return nil
}
