package shannon_fano

import (
	"fmt"
	"iden69/archiver/lib/compression/table"
	"math"
	"sort"
	"strings"
)

type Generator struct{}

func NewGenerator() Generator {
	return Generator{}
}

func (g Generator) NewTable(text string) table.EncodingTable {
	stat := newCharStat(text)

	t := build(stat)

	return t.Export()
}

func (et encodingTable) Export() map[rune]string {
	res := make(map[rune]string)

	for k, v := range et {
		byteStr := fmt.Sprintf("%b", v.Bits)

		if lenDiff := v.Size - len(byteStr); lenDiff > 0 {
			byteStr = strings.Repeat("0", lenDiff) + byteStr
		}

		res[k] = byteStr
	}

	return res
}

type code struct {
	Char     rune
	Quantity int
	Bits     uint32
	Size     int
}

type encodingTable map[rune]code

func build(stat charStat) encodingTable {
	codes := make([]code, 0, len(stat))

	for ch, qty := range stat {
		codes = append(codes, code{
			Char:     ch,
			Quantity: qty,
		})
	}

	sort.Slice(codes, func(i, j int) bool {
		if codes[i].Quantity != codes[j].Quantity {
			return codes[i].Quantity > codes[j].Quantity
		}

		return codes[i].Char < codes[j].Char
	})

	assignCodes(codes)

	res := make(encodingTable)

	for _, code := range codes {
		res[code.Char] = code
	}

	return res
}

func assignCodes(codes []code) {
	if len(codes) < 2 {
		return
	}

	divider := bestDividerPosition(codes)

	for i := 0; i < len(codes); i++ {
		codes[i].Bits <<= 1
		codes[i].Size++

		if i >= divider {
			codes[i].Bits |= 1
		}
	}

	assignCodes(codes[:divider])
	assignCodes(codes[divider:])
}

func bestDividerPosition(codes []code) int {
	total := 0
	left := 0
	bestPosition := 0
	prevDiff := math.MaxInt

	for _, code := range codes {
		total += code.Quantity
	}

	for i := 0; i < len(codes)-1; i++ {
		left += codes[i].Quantity
		right := total - left

		diff := abs(right - left)
		if diff >= prevDiff {
			break
		}

		prevDiff = diff
		bestPosition = i + 1
	}

	return bestPosition
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

type charStat map[rune]int

func newCharStat(text string) charStat {
	res := make(charStat)

	for _, ch := range text {
		res[ch]++
	}

	return res
}
