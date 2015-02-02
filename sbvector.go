package sbvector

import (
	"github.com/hideo55/go-popcount"
)

// SBVector is succinct bit vector type.
type SBVector struct {
	blocks       []uint64
	ranks        []rankindex
	select1Table []uint64
	select0Table []uint64
	numOf1s      uint64
	size         uint64
}

const (
	mask55     uint64 = 0x5555555555555555
	mask33     uint64 = 0x3333333333333333
	mask0F     uint64 = 0x0F0F0F0F0F0F0F0F
	mask01     uint64 = 0x0101010101010101
	mask80     uint64 = 0x8080808080808080
	sBlockSize uint64 = 64
	lBlockSize uint64 = 512
	blockRate  uint64 = 8
	notFound   uint64 = 0xFFFFFFFFFFFFFFFF
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

type sbvecMethods interface {
	set(i uint64, val bool)
	get(i uint64) uint64
	rank(i uint64) uint64
	select1(x uint64) uint64
	getSize(b bool) uint64
	build(enableFasterSelect1 bool, enableFasterSelect0 bool)
}

func (vec *SBVector) get(i uint64) bool {
	return (vec.blocks[i/sBlockSize] & (1 << (i % sBlockSize))) != 0
}

func (vec *SBVector) set(i uint64, val bool) {
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
}

func (vec *SBVector) build(enableFasterSelect1 bool, enableFasterSelect0 bool) {
	var blockNum = uint64(len(vec.blocks))
	var numOf1s = lBlockSize
	var numOf0s = lBlockSize
	vec.numOf1s = 0

	clearSlice(vec.select1Table)
	clearSlice(vec.select0Table)
	var tmpRank = make([]rankindex, 1)
	copy(tmpRank, vec.ranks)
	var rankTableSize = (blockNum * sBlockSize) / lBlockSize  + 1
	if ((blockNum * sBlockSize) % lBlockSize) != 0 {
		rankTableSize++
	}

	vec.ranks = make([]rankindex, rankTableSize)

	for i, x := range vec.blocks {
		var rankID = uint64(i) / blockRate
		var rank  = &vec.ranks[rankID]
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
			rank.setRel3(2/*vec.numOf1s - rank.abs()*/)
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

func (vec *SBVector) rank(i uint64) uint64 {
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
	return offset
}

func (vec *SBVector) select1(x uint64) uint64 {
	var vecSize = vec.getSize(true)
	if vecSize <= x {
		return notFound
	}

	var begin uint64 = 0
	var end uint64 = 0

	if len(vec.select1Table) == 0 {
		begin = 0
		end = uint64(len(vec.ranks))
	} else {
		var selectID = x / lBlockSize
		if x&lBlockSize == 0 {
			return vec.select1Table[selectID]
		}
		begin = vec.select1Table[selectID] / lBlockSize
		end = (vec.select1Table[selectID+1] + lBlockSize - 1) / lBlockSize
	}

	if (begin + 10) >= end {
		for x >= vec.ranks[begin + 1].abs() {
			begin++
		}
	} else {
		for  (begin + 1) < end {
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
            blockID += 4;
            x -= rank.rel4();
        } else {
            blockID += 5;
            x -= rank.rel5();
        }
    } else if x < rank.rel7() {
        blockID += 6;
        x -= rank.rel6();
    } else {
        blockID += 7;
        x -= rank.rel7();
    }
	return select64(vec.blocks[blockID], x, blockID * sBlockSize)
}

func (vec *SBVector) getSize(b bool) uint64{
	if b {
		return vec.numOf1s
	} else {
		return vec.size - vec.numOf1s
	}
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

func clearSlice(s []uint64) {
	var c = make([]uint64, 0)
	copy(c, s)
}
