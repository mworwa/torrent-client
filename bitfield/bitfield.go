package bitfield

import "math"

type Bitfield []byte

func New(payload []byte) Bitfield {
	return Bitfield(payload[1:])
}

func NewEmpty(lenght int, pieceLenght int) Bitfield {

	pieceCount := int(math.Ceil(float64(lenght) / float64(pieceLenght)))
	numBytes := (pieceCount + 7) / 8
	return make(Bitfield, numBytes)
}

func (b Bitfield) hasPiece(index int) bool {
	byteIndex := index / 8
	bitOffset := index % 8

	return b[byteIndex]>>uint(7-bitOffset)&1 != 0
}
