package processors

import (
	"errors"
	"flag"
	"fmt"
)

const (
	outFilenamePostfix = ".out"
)

const (
	modeEncode  = "encode"
	modeDecode  = "decode"
	modeEntropy = "entropy"
)

type flags struct {
	mode           string
	InputFilename  string
	OutputFilename string
}

type FlagProcessor struct {
	flags
}

func NewFlagProcessor() *FlagProcessor {
	FlagProcessor := new(FlagProcessor)

	return FlagProcessor
}

func (processor *FlagProcessor) ProcessInput() error {
	// TODO: сделать без указания флагов
	modeFlag := flag.String("mode", modeEncode, "specify work mode")
	InputFilenameFlag := flag.String("i", "", "specify input filename")
	OutputFilenameFlag := flag.String("o", "", "specify output filename")

	flag.Parse()

	processor.mode = *modeFlag
	processor.InputFilename = *InputFilenameFlag
	processor.OutputFilename = *OutputFilenameFlag

	err := processor.validateInput()

	if err != nil {
		return err
	}

	if processor.OutputFilename == "" {
		processor.OutputFilename = processor.InputFilename + outFilenamePostfix
	}

	return nil
}

func (processor *FlagProcessor) IsEncodeMode() bool {
	return processor.mode == modeEncode
}

func (processor *FlagProcessor) IsDecodeMode() bool {
	return processor.mode == modeDecode
}

func (processor *FlagProcessor) IsEntropyMode() bool {
	return processor.mode == modeEntropy
}

func (processor *FlagProcessor) validateInput() error {
	mode := processor.mode
	InputFilename := processor.InputFilename

	if mode != modeDecode && mode != modeEncode && mode != modeEntropy {
		return errors.New(fmt.Sprintf("error: invalid mode - %s\n", mode))
	}

	if InputFilename == "" {
		return errors.New(fmt.Sprintln("error: input filename not specified"))
	}

	return nil
}
