package main

import (
	"flag"
	"fmt"
	"huffman/src/decoder"
	"huffman/src/encoder"
	"os"
)

const (
	ModeEncode = "encode"
	ModeDecode = "decode"
)

func main() {
	modeFlag := flag.String("mode", ModeEncode, "specify work mode")
	// TODO: сделать без указания флагов
	inputFilenameFlag := flag.String("i", "", "specify input filename")
	outputFilenameFlag := flag.String("o", "", "specify output filename")
	flag.Parse()

	mode := *modeFlag
	inputFilename := *inputFilenameFlag
	outputFilename := *outputFilenameFlag

	// TODO: перенести из main, сделать на error
	if mode == ModeDecode {
		decoder.Decode()
	} else if mode == ModeEncode {
		encoder.Encode()
	} else {
		fmt.Printf("error: invalid mode - %s\n", mode)
		os.Exit(1)
	}

	if inputFilename == "" {
		fmt.Println("error: input filename not specified")
		os.Exit(2)
	}

	if outputFilename == "" {
		outputFilename = inputFilename + ".out" // TODO: убрать хардкод
	}

}
