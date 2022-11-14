package entropy

import (
	"bufio"
	"io"
	"math"
)

type Entropy struct {
	wordsTotal   uint64
	wordCounters map[byte]*counter
}

type counter struct {
	count               uint64
	countPrecedingWords uint64
	countPrecedingPairs uint64
	precedingWords      map[byte]uint64
	precedingPairs      map[uint16]uint64
}

func newCounter() *counter {
	counter := new(counter)
	counter.precedingWords = make(map[byte]uint64)
	counter.precedingPairs = make(map[uint16]uint64)

	return counter
}

func NewEntropy() *Entropy {
	entropy := &Entropy{}

	entropy.wordCounters = make(map[byte]*counter)

	return entropy
}

func (entropy *Entropy) Init(reader *bufio.Reader) {
	var prevByte byte
	var prevSecondByte byte
	prevByteSet := false
	prevSecondByteSet := false

	for {
		newByte, err := reader.ReadByte()

		if err == io.EOF {
			break
		}

		entropy.wordsTotal++

		if counter, ok := entropy.wordCounters[newByte]; ok {
			counter.count++
		} else {
			entropy.wordCounters[newByte] = newCounter()
			entropy.wordCounters[newByte].count++
		}

		if prevByteSet {
			if prevSecondByteSet {
				entropy.wordCounters[newByte].precedingPairs[(uint16(prevSecondByte)<<8)|uint16(prevByte)]++
				entropy.wordCounters[newByte].countPrecedingPairs++
			}
			entropy.wordCounters[newByte].precedingWords[prevByte]++
			entropy.wordCounters[newByte].countPrecedingWords++
			prevSecondByteSet = true
			prevSecondByte = prevByte
		}

		prevByteSet = true
		prevByte = newByte
	}
}

func (entropy *Entropy) CalculateEntropy() float64 { //ok
	var E float64

	for _, counter := range entropy.wordCounters {
		p := float64(counter.count) / float64(entropy.wordsTotal)

		E += -(p * math.Log2(p))
	}

	return E
}

func (entropy *Entropy) CalculateBlockEntropy(E float64, E_XY float64, blockSize uint64) float64 {
	return (E + E_XY*(float64(blockSize)-1)) / float64(blockSize)
}

func (entropy *Entropy) CalculateConditionalEntropy() float64 {
	var E_XY float64

	for _, counter := range entropy.wordCounters {
		p := float64(counter.count) / float64(entropy.wordsTotal)
		total := float64(counter.countPrecedingWords)
		var sum float64

		for _, count := range counter.precedingWords {
			pCond := float64(count) / total
			log := math.Log2(pCond)
			sum += pCond * log
		}

		E_XY += -p * sum
	}

	return E_XY
}

func (entropy *Entropy) CalculateDoubleConditionalEntropy() float64 {
	var E_XYY float64

	for _, counter := range entropy.wordCounters {
		p := float64(counter.count) / float64(entropy.wordsTotal)
		total := float64(counter.countPrecedingPairs)
		var sum float64

		for _, count := range counter.precedingPairs {
			pCond := float64(count) / total
			log := math.Log2(pCond)
			sum += pCond * log
		}

		E_XYY += -p * sum
	}

	return E_XYY
}
