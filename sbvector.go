/*
Package sbvector is implementation of succinct bit vector for Go.

Synopsis
	import (
		"github.com/hideo55/go-sbvector"
	)

	func example() {
		vec, err := sbvector.NewVector()
		if err != nil {
			// error handling
		}

		vec.Set(10, true)

		...

		vec.Build(true, true)

		pos, err := vec.Select1(0)
		if err != nil {
			// error handling
		}
	}
*/
package sbvector

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"unsafe"

	"github.com/hideo55/go-popcount"
)

// BitVectorData holds impormation about bit vector.
type BitVectorData struct {
	blocks       []uint64
	ranks        []rankIndex
	select1Table []uint64
	select0Table []uint64
	numOf1s      uint64
	size         uint64
}

// SuccinctBitVector is interface of succinct bit vector.
type SuccinctBitVector interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Set(i uint64, val bool)
	PushBack(b bool)
	PushBackBits(x uint64, length uint64)
	Get(i uint64) (bool, error)
	GetBits(pos uint64, length uint64) (uint64, error)
	Rank1(i uint64) (uint64, error)
	Rank0(i uint64) (uint64, error)
	Rank(i uint64, b bool) (uint64, error)
	Select1(x uint64) (uint64, error)
	Select0(x uint64) (uint64, error)
	Select(x uint64, b bool) (uint64, error)
	Size() uint64
	NumOfBits(b bool) uint64
	Build(enableFasterSelect1 bool, enableFasterSelect0 bool)
}

const (
	mask55      uint64 = 0x5555555555555555
	mask33      uint64 = 0x3333333333333333
	mask0F      uint64 = 0x0F0F0F0F0F0F0F0F
	mask01      uint64 = 0x0101010101010101
	mask80      uint64 = 0x8080808080808080
	sBlockSize  uint64 = 64
	lBlockSize  uint64 = 512
	blockRate   uint64 = 8
	minimumSize uint64 = 40
	// NotFound indicates `value not found`
	NotFound uint64 = 0xFFFFFFFFFFFFFFFF
)

var selectTable = [8][256]uint8{
	[256]uint8{7, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 5, 0, 1, 0, 2, 0, 1,
		0, 3, 0, 1, 0, 2, 0, 1, 0, 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 6, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0,
		2, 0, 1, 0, 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 4,
		0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 7, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 4, 0, 1, 0, 2, 0,
		1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1,
		0, 2, 0, 1, 0, 6, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0,
		5, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0, 4, 0, 1, 0, 2, 0, 1, 0, 3, 0, 1, 0, 2, 0, 1, 0},
	[256]uint8{7, 7, 7, 1, 7, 2, 2, 1, 7, 3, 3, 1, 3, 2, 2, 1, 7, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1, 7, 5, 5, 1, 5, 2, 2,
		1, 5, 3, 3, 1, 3, 2, 2, 1, 5, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1, 7, 6, 6, 1, 6, 2, 2, 1, 6, 3, 3, 1,
		3, 2, 2, 1, 6, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1, 6, 5, 5, 1, 5, 2, 2, 1, 5, 3, 3, 1, 3, 2, 2, 1, 5,
		4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1, 7, 7, 7, 1, 7, 2, 2, 1, 7, 3, 3, 1, 3, 2, 2, 1, 7, 4, 4, 1, 4, 2,
		2, 1, 4, 3, 3, 1, 3, 2, 2, 1, 7, 5, 5, 1, 5, 2, 2, 1, 5, 3, 3, 1, 3, 2, 2, 1, 5, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3,
		1, 3, 2, 2, 1, 7, 6, 6, 1, 6, 2, 2, 1, 6, 3, 3, 1, 3, 2, 2, 1, 6, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1,
		6, 5, 5, 1, 5, 2, 2, 1, 5, 3, 3, 1, 3, 2, 2, 1, 5, 4, 4, 1, 4, 2, 2, 1, 4, 3, 3, 1, 3, 2, 2, 1},
	[256]uint8{7, 7, 7, 7, 7, 7, 7, 2, 7, 7, 7, 3, 7, 3, 3, 2, 7, 7, 7, 4, 7, 4, 4, 2, 7, 4, 4, 3, 4, 3, 3, 2, 7, 7, 7, 5, 7, 5, 5,
		2, 7, 5, 5, 3, 5, 3, 3, 2, 7, 5, 5, 4, 5, 4, 4, 2, 5, 4, 4, 3, 4, 3, 3, 2, 7, 7, 7, 6, 7, 6, 6, 2, 7, 6, 6, 3,
		6, 3, 3, 2, 7, 6, 6, 4, 6, 4, 4, 2, 6, 4, 4, 3, 4, 3, 3, 2, 7, 6, 6, 5, 6, 5, 5, 2, 6, 5, 5, 3, 5, 3, 3, 2, 6,
		5, 5, 4, 5, 4, 4, 2, 5, 4, 4, 3, 4, 3, 3, 2, 7, 7, 7, 7, 7, 7, 7, 2, 7, 7, 7, 3, 7, 3, 3, 2, 7, 7, 7, 4, 7, 4,
		4, 2, 7, 4, 4, 3, 4, 3, 3, 2, 7, 7, 7, 5, 7, 5, 5, 2, 7, 5, 5, 3, 5, 3, 3, 2, 7, 5, 5, 4, 5, 4, 4, 2, 5, 4, 4,
		3, 4, 3, 3, 2, 7, 7, 7, 6, 7, 6, 6, 2, 7, 6, 6, 3, 6, 3, 3, 2, 7, 6, 6, 4, 6, 4, 4, 2, 6, 4, 4, 3, 4, 3, 3, 2,
		7, 6, 6, 5, 6, 5, 5, 2, 6, 5, 5, 3, 5, 3, 3, 2, 6, 5, 5, 4, 5, 4, 4, 2, 5, 4, 4, 3, 4, 3, 3, 2},
	[256]uint8{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 3, 7, 7, 7, 7, 7, 7, 7, 4, 7, 7, 7, 4, 7, 4, 4, 3, 7, 7, 7, 7, 7, 7, 7,
		5, 7, 7, 7, 5, 7, 5, 5, 3, 7, 7, 7, 5, 7, 5, 5, 4, 7, 5, 5, 4, 5, 4, 4, 3, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6,
		7, 6, 6, 3, 7, 7, 7, 6, 7, 6, 6, 4, 7, 6, 6, 4, 6, 4, 4, 3, 7, 7, 7, 6, 7, 6, 6, 5, 7, 6, 6, 5, 6, 5, 5, 3, 7,
		6, 6, 5, 6, 5, 5, 4, 6, 5, 5, 4, 5, 4, 4, 3, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 3, 7, 7, 7, 7, 7, 7,
		7, 4, 7, 7, 7, 4, 7, 4, 4, 3, 7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 5, 7, 5, 5, 3, 7, 7, 7, 5, 7, 5, 5, 4, 7, 5, 5,
		4, 5, 4, 4, 3, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 3, 7, 7, 7, 6, 7, 6, 6, 4, 7, 6, 6, 4, 6, 4, 4, 3,
		7, 7, 7, 6, 7, 6, 6, 5, 7, 6, 6, 5, 6, 5, 5, 3, 7, 6, 6, 5, 6, 5, 5, 4, 6, 5, 5, 4, 5, 4, 4, 3},
	[256]uint8{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 4, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 5, 7, 5, 5, 4, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 6, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 4, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 5, 7,
		7, 7, 6, 7, 6, 6, 5, 7, 6, 6, 5, 6, 5, 5, 4, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 4, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7,
		5, 7, 5, 5, 4, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 4,
		7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 5, 7, 7, 7, 6, 7, 6, 6, 5, 7, 6, 6, 5, 6, 5, 5, 4},
	[256]uint8{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 5, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6, 7,
		7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 5, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 5, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 6, 7, 6, 6, 5},
	[256]uint8{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6},
	[256]uint8{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7,
		7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 6},
}

// NewVector returns a new succinct bit vector
func NewVector() (SuccinctBitVector, error) {
	vec := new(BitVectorData)
	return vec, nil
}

// Get returns value from bit vector by index.
func (vec *BitVectorData) Get(i uint64) (bool, error) {
	if i > vec.size {
		return false, errors.New("Out of Bounds")
	}
	return (vec.blocks[i/sBlockSize] & (1 << (i % sBlockSize))) != 0, nil
}

//GetBits returns bits from bit vector.
func (vec *BitVectorData) GetBits(pos uint64, length uint64) (uint64, error) {
	if (pos + length) > vec.size {
		return NotFound, errors.New("Out of Bounds")
	}
	var blockIdx1 = pos / sBlockSize
	var blockOffset1 = pos % sBlockSize
	if (blockOffset1 + length) <= sBlockSize {
		return mask(vec.blocks[blockIdx1]>>blockOffset1, length), nil
	}
	var blockIdx2 = (pos + length - 1) / sBlockSize
	return mask((vec.blocks[blockIdx1]>>blockOffset1)+(vec.blocks[blockIdx2]<<(sBlockSize-blockOffset1)), length), nil
}

// Set value to bit vector by index.
func (vec *BitVectorData) Set(i uint64, val bool) {
	var blockID uint64
	var r uint8
	blockID = i / sBlockSize
	r = uint8(i % sBlockSize)
	if uint64(len(vec.blocks)) <= blockID {
		newSlice := make([]uint64, (blockID - uint64(len(vec.blocks)) + 1))
		vec.blocks = append(vec.blocks, newSlice...)
	}
	var m uint64
	m = 0x1 << r
	if val {
		vec.blocks[blockID] |= m
	} else {
		vec.blocks[blockID] &= ^m
	}
	if vec.size <= i {
		vec.size = i + 1
	}
}

// PushBack add bit to the bit vector
func (vec *BitVectorData) PushBack(b bool) {
	if (vec.size / sBlockSize) >= uint64(len(vec.blocks)) {
		vec.blocks = append(vec.blocks, uint64(0))
	}
	var blockID = vec.size / sBlockSize
	var r = vec.size % sBlockSize
	var m = uint64(1) << r
	if b {
		vec.blocks[blockID] |= m
	} else {
		vec.blocks[blockID] &= ^m
	}
	vec.size++
}

// PushBackBits add bits to the bit vector
func (vec *BitVectorData) PushBackBits(x uint64, length uint64) {
	var offset = vec.size % sBlockSize
	if (vec.size+length-1)/sBlockSize >= uint64(len(vec.blocks)) {
		vec.blocks = append(vec.blocks, uint64(0))
	}
	var blockID = vec.size / sBlockSize
	vec.blocks[blockID] |= (x << offset)
	if (offset + length - 1) >= sBlockSize {
		vec.blocks[blockID+1] |= (x >> (sBlockSize - offset))
	}
	vec.size += length
}

// Build creates indexes for succinct bit vector(rank index, ...).
// If `enableFasterSelect1` is true, creates index for select1 make faster.
// If `enableFasterSelect0` is true, creates index for select0 make faster.
func (vec *BitVectorData) Build(enableFasterSelect1 bool, enableFasterSelect0 bool) {
	var blockNum = uint64(len(vec.blocks))
	var numOf1s = lBlockSize
	var numOf0s = lBlockSize
	vec.numOf1s = 0

	clearSlice(vec.select1Table)
	clearSlice(vec.select0Table)
	var tmpRank = make([]rankIndex, 1)
	copy(tmpRank, vec.ranks)
	var rankTableSize = (blockNum*sBlockSize)/lBlockSize + 1
	if ((blockNum * sBlockSize) % lBlockSize) != 0 {
		rankTableSize++
	}

	vec.ranks = make([]rankIndex, rankTableSize)

	for i, x := range vec.blocks {
		var rankID = uint64(i) / blockRate
		var rank = &vec.ranks[rankID]
		switch i % 8 {
		case 0:
			rank.setAbs(vec.numOf1s)
		case 1:
			rank.setRel1(vec.numOf1s - rank.abs())
		case 2:
			rank.setRel2(vec.numOf1s - rank.abs())
		case 3:
			rank.setRel3(vec.numOf1s - rank.abs())
		case 4:
			rank.setRel4(vec.numOf1s - rank.abs())
		case 5:
			rank.setRel5(vec.numOf1s - rank.abs())
		case 6:
			rank.setRel6(vec.numOf1s - rank.abs())
		case 7:
			rank.setRel7(vec.numOf1s - rank.abs())
		}

		var count1s = popcount.Count(x)
		if enableFasterSelect1 && (numOf1s+count1s > lBlockSize) {
			var diff = lBlockSize - numOf1s
			var pos = select64(x, diff, 0)
			vec.select1Table = append(vec.select1Table, uint64(i)*sBlockSize+pos)
			numOf1s -= lBlockSize
		}

		var count0s = sBlockSize - count1s
		if enableFasterSelect0 && (numOf0s+count0s > lBlockSize) {
			var diff = lBlockSize - numOf0s
			var pos = select64(^x, diff, 0)
			vec.select0Table = append(vec.select0Table, uint64(i)*sBlockSize+pos)
			numOf0s -= lBlockSize
		}

		numOf1s += count1s
		numOf0s += count0s
		vec.numOf1s += count1s
	}

	if (blockNum % blockRate) != 0 {
		var rankID = (blockNum - 1) / blockRate
		var rank = &vec.ranks[rankID]
		switch (blockNum - 1) % blockRate {
		case 0:
			rank.setRel1(vec.numOf1s - rank.abs())
			fallthrough
		case 1:
			rank.setRel2(vec.numOf1s - rank.abs())
			fallthrough
		case 2:
			rank.setRel3(2 /*vec.numOf1s - rank.abs()*/)
			fallthrough
		case 3:
			rank.setRel4(vec.numOf1s - rank.abs())
			fallthrough
		case 4:
			rank.setRel5(vec.numOf1s - rank.abs())
			fallthrough
		case 5:
			rank.setRel6(vec.numOf1s - rank.abs())
			fallthrough
		case 6:
			rank.setRel7(vec.numOf1s - rank.abs())
		}
	}

	vec.ranks[len(vec.ranks)-1].setAbs(vec.numOf1s)

	if enableFasterSelect1 {
		vec.select1Table = append(vec.select1Table, vec.size)
	}
	if enableFasterSelect0 {
		vec.select0Table = append(vec.select0Table, vec.size)
	}
}

// Rank1 returns number of the bits equal to `1` up to positin `i`
func (vec *BitVectorData) Rank1(i uint64) (uint64, error) {
	if i > vec.size {
		return NotFound, errors.New("Out of Bounds")
	}
	var rankID = i / lBlockSize
	var blockID = i / sBlockSize
	var r = i % sBlockSize

	var rank = vec.ranks[rankID]
	var offset = rank.abs()
	switch blockID % blockRate {
	case 1:
		offset += rank.rel1()
	case 2:
		offset += rank.rel2()
	case 3:
		offset += rank.rel3()
	case 4:
		offset += rank.rel4()
	case 5:
		offset += rank.rel5()
	case 6:
		offset += rank.rel6()
	case 7:
		offset += rank.rel7()
	default:
	}
	offset += popcount.Count(vec.blocks[blockID] & ((1 << r) - 1))
	return offset, nil
}

// Rank0 returns number of the bits equal to `0` up to positin `i`
func (vec *BitVectorData) Rank0(i uint64) (uint64, error) {
	rank, err := vec.Rank1(i)
	if err != nil {
		return rank, err
	}
	return i - rank, nil
}

// Rank returns number of the bits equal to `b` up to position `i`
func (vec *BitVectorData) Rank(i uint64, b bool) (uint64, error) {
	if b {
		return vec.Rank1(i)
	}
	return vec.Rank0(i)
}

// Select1 returns the position of the x-th occurence of 1
func (vec *BitVectorData) Select1(x uint64) (uint64, error) {
	var vecSize = vec.NumOfBits(true)
	if vecSize <= x {
		return NotFound, errors.New("Out of Bounds")
	}

	var begin uint64
	var end uint64

	if len(vec.select1Table) == 0 {
		begin = 0
		end = uint64(len(vec.ranks))
	} else {
		var selectID = x / lBlockSize
		if x%lBlockSize == 0 {
			return vec.select1Table[selectID], nil
		}
		begin = vec.select1Table[selectID] / lBlockSize
		end = (vec.select1Table[selectID+1] + lBlockSize - 1) / lBlockSize
	}
	if (begin + 10) >= end {
		for x >= vec.ranks[begin+1].abs() {
			begin++
		}
	} else {
		for (begin + 1) < end {
			var pivot = (begin + end) / 2
			if x < vec.ranks[pivot].abs() {
				end = pivot
			} else {
				begin = pivot
			}
		}
	}
	var rankID = begin
	var rank = &vec.ranks[rankID]
	var rankOffset = rank.abs()
	x -= rankOffset
	var blockID = rankID * blockRate
	if x < rank.rel4() {
		if x < rank.rel2() {
			if x >= rank.rel1() {
				blockID++
				x -= rank.rel1()
			}
		} else if x < rank.rel3() {
			blockID += 2
			x -= rank.rel2()
		} else {
			blockID += 3
			x -= rank.rel3()
		}
	} else if x < rank.rel6() {
		if x < rank.rel5() {
			blockID += 4
			x -= rank.rel4()
		} else {
			blockID += 5
			x -= rank.rel5()
		}
	} else if x < rank.rel7() {
		blockID += 6
		x -= rank.rel6()
	} else {
		blockID += 7
		x -= rank.rel7()
	}
	return select64(vec.blocks[blockID], x, blockID*sBlockSize), nil
}

// Select0 returns the position of the x-th occurence of 0
func (vec *BitVectorData) Select0(x uint64) (uint64, error) {
	var vecSize = vec.NumOfBits(false)
	if vecSize <= x {
		return NotFound, errors.New("Out of Bounds")
	}

	var begin uint64
	var end uint64

	if len(vec.select0Table) == 0 {
		begin = 0
		end = uint64(len(vec.ranks))
	} else {
		var selectID = x / lBlockSize
		if x%lBlockSize == 0 {
			return vec.select0Table[selectID], nil
		}
		begin = vec.select0Table[selectID] / lBlockSize
		end = (vec.select0Table[selectID+1] + lBlockSize - 1) / lBlockSize
	}

	if (begin + 10) >= end {
		for x >= ((begin+1)*lBlockSize)-vec.ranks[begin+1].abs() {
			begin++
		}
	} else {
		for (begin + 1) < end {
			var pivot = (begin + end) / 2
			if x < (pivot*lBlockSize)-vec.ranks[pivot].abs() {
				end = pivot
			} else {
				begin = pivot
			}
		}
	}
	var rankID = begin
	var rank = &vec.ranks[rankID]
	var rankOffset = (rankID * lBlockSize) - rank.abs()
	x -= rankOffset
	var blockID = rankID * blockRate
	if x < uint64(256)-rank.rel4() {
		if x < uint64(128)-rank.rel2() {
			if x >= uint64(64)-rank.rel1() {
				blockID++
				x -= uint64(64) - rank.rel1()
			}
		} else if x < uint64(192)-rank.rel3() {
			blockID += 2
			x -= uint64(128) - rank.rel2()
		} else {
			blockID += 3
			x -= uint64(192) - rank.rel3()
		}
	} else if x < uint64(384)-rank.rel6() {
		if x < uint64(320)-rank.rel5() {
			blockID += 4
			x -= uint64(256) - rank.rel4()
		} else {
			blockID += 5
			x -= uint64(320) - rank.rel5()
		}
	} else if x < uint64(448)-rank.rel7() {
		blockID += 6
		x -= uint64(384) - rank.rel6()
	} else {
		blockID += 7
		x -= uint64(448) - rank.rel7()
	}
	return select64(^vec.blocks[blockID], x, blockID*sBlockSize), nil
}

// Select returns the position of the x-th occurrence of `b`
func (vec *BitVectorData) Select(x uint64, b bool) (uint64, error) {
	if b {
		return vec.Select1(x)
	}
	return vec.Select0(x)
}

// Size returns size of bit vector
func (vec *BitVectorData) Size() uint64 {
	return vec.size
}

// NumOfBits returns number of bits that matches with argument in the bit vector.
func (vec *BitVectorData) NumOfBits(b bool) uint64 {
	if b {
		return vec.numOf1s
	}
	return vec.size - vec.numOf1s
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (vec *BitVectorData) MarshalBinary() ([]byte, error) {
	buffer := new(bytes.Buffer)

	blockNum := uint32(len(vec.blocks))
	rankIndexSize := uint32(len(vec.ranks))
	select1TableSize := uint32(len(vec.select1Table))
	select0TableSize := uint32(len(vec.select0Table))

	var tmpUint32 uint32
	var tmpRankIndex rankIndex
	sizeOfUint32 := uint64(unsafe.Sizeof(tmpUint32))
	sizeOfUint64 := uint64(unsafe.Sizeof(vec.size))
	sizeOfRI := uint64(unsafe.Sizeof(tmpRankIndex))

	var serializedSize uint64
	serializedSize = uint64(blockNum) * sizeOfUint64
	serializedSize += uint64(rankIndexSize) * sizeOfRI
	serializedSize += uint64(select1TableSize) * sizeOfUint64
	serializedSize += uint64(select0TableSize) * sizeOfUint64
	serializedSize += sizeOfUint64 * 3 /* Sizeof(serializedSize) + Sizeof(vec.size) + Sizeof(vec.numOf1s) */
	serializedSize += sizeOfUint32 * 4 /* Sizeof(blockNum) + Sizeof(rankIndexSize) + Sizeof(select1TableSize) + Sizeof(select0TableSize) */
	binary.Write(buffer, binary.LittleEndian, &serializedSize)
	binary.Write(buffer, binary.LittleEndian, &vec.size)
	binary.Write(buffer, binary.LittleEndian, &vec.numOf1s)

	binary.Write(buffer, binary.LittleEndian, &blockNum)
	for _, block := range vec.blocks {
		binary.Write(buffer, binary.LittleEndian, &block)
	}

	binary.Write(buffer, binary.LittleEndian, &rankIndexSize)
	for _, ri := range vec.ranks {
		buf, err := ri.MarshalBinary()
		if err != nil {
			return make([]byte, 0), err
		}
		binary.Write(buffer, binary.LittleEndian, buf)
	}
	binary.Write(buffer, binary.LittleEndian, &select1TableSize)
	for _, s1 := range vec.select1Table {
		binary.Write(buffer, binary.LittleEndian, &s1)
	}
	binary.Write(buffer, binary.LittleEndian, &select0TableSize)
	for _, s0 := range vec.select0Table {
		binary.Write(buffer, binary.LittleEndian, &s0)
	}
	return buffer.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (vec *BitVectorData) UnmarshalBinary(data []byte) error {
	buf := data
	if uint64(len(data)) < minimumSize {
		return errors.New("no data")
	}
	var offset int
	offset = 0
	buf = data[offset : offset+8]
	offset += 8
	dataSize := binary.LittleEndian.Uint64(buf)
	if uint64(len(data)) != dataSize {
		return errors.New("Invalid Size")
	}

	buf = data[offset : offset+8]
	offset += 8
	vec.size = binary.LittleEndian.Uint64(buf)

	buf = data[offset : offset+8]
	offset += 8
	vec.numOf1s = binary.LittleEndian.Uint64(buf)

	buf = data[offset : offset+4]
	offset += 4
	blockNum := binary.LittleEndian.Uint32(buf)
	vec.blocks = make([]uint64, blockNum)
	for i := uint32(0); i < blockNum; i++ {
		buf = data[offset : offset+8]
		offset += 8
		vec.blocks[i] = binary.LittleEndian.Uint64(buf)
	}

	buf = data[offset : offset+4]
	offset += 4
	rankIndexSize := binary.LittleEndian.Uint32(buf)
	vec.ranks = make([]rankIndex, rankIndexSize)
	sizeOfRI := unsafe.Sizeof(vec.ranks[0])
	for i := uint32(0); i < rankIndexSize; i++ {
		buf = data[offset : offset+int(sizeOfRI)]
		offset += int(sizeOfRI)
		vec.ranks[i].UnmarshalBinary(buf)
	}

	buf = data[offset : offset+4]
	offset += 4
	select1TableSize := binary.LittleEndian.Uint32(buf)
	vec.select1Table = make([]uint64, select1TableSize)
	for i := uint32(0); i < select1TableSize; i++ {
		buf = data[offset : offset+8]
		offset += 8
		vec.select1Table[i] = binary.LittleEndian.Uint64(buf)
	}

	buf = data[offset : offset+4]
	offset += 4
	select0TableSize := binary.LittleEndian.Uint32(buf)
	vec.select0Table = make([]uint64, select0TableSize)
	for i := uint32(0); i < select0TableSize; i++ {
		buf = data[offset : offset+8]
		offset += 8
		vec.select0Table[i] = binary.LittleEndian.Uint64(buf)
	}

	return nil
}

func countTrailingZeros(x uint64) uint8 {
	return uint8(popcount.Count((x & (-x)) - 1))
}

func select64(block uint64, i uint64, base uint64) uint64 {
	var counts uint64
	counts = block - ((block >> 1) & mask55)
	counts = (counts & mask33) + ((counts >> 2) & mask33)
	counts = (counts + (counts >> 4)) & mask0F
	counts *= mask01

	var x = (counts | mask80) - ((i + 1) * mask01)
	var tzLen = countTrailingZeros((x & mask80) >> 7)
	base += uint64(tzLen)
	block >>= tzLen
	i -= ((counts << 8) >> tzLen) & 0xFF
	return base + uint64(selectTable[i][block&0xFF])
}

func mask(x uint64, pos uint64) uint64 {
	return x & ((uint64(1) << pos) - 1)
}

func clearSlice(s []uint64) {
	var c = make([]uint64, 0)
	copy(c, s)
}
