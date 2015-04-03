package sbvector

import (
	"encoding/binary"
	"testing"
)

type BitTestCase struct {
	pos uint64
	bit bool
}

type RankTestCase struct {
	pos  uint64
	rank uint64
}

type SelectTestCase struct {
	index uint64
	pos   uint64
}

var (
	bitCases = []BitTestCase{
		{0, true},
		{1, true},
		{63, true},
		{64, true},
		{65, true},
		{127, true},
		{128, true},
		{191, true},
		{192, true},
		{255, true},
		{256, true},
		{319, true},
		{320, true},
		{383, true},
		{384, true},
		{447, true},
		{448, true},
		{511, true},
		{512, true},
		{1024, false},
		{6000, true},
	}

	rankCases = []RankTestCase{
		{1, 1},
		{2, 2},
		{3, 2},
		{63, 2},
		{64, 3},
		{65, 4},
		{66, 5},
		{128, 6},
		{129, 7},
		{192, 8},
		{193, 9},
		{255, 9},
		{256, 10},
		{319, 11},
		{320, 12},
		{383, 13},
		{384, 14},
		{447, 15},
		{448, 16},
		{511, 17},
		{512, 18},
		{513, 19},
		{1024, 19},
		{6001, 20},
	}

	select1Cases = []SelectTestCase{
		{0, 0},
		{1, 1},
		{2, 63},
		{3, 64},
		{4, 65},
		{5, 127},
		{6, 128},
		{7, 191},
		{8, 192},
		{9, 255},
		{10, 256},
		{11, 319},
		{12, 320},
		{13, 383},
		{14, 384},
		{15, 447},
		{16, 448},
		{17, 511},
		{18, 512},
		{19, 6000},
	}

	select0Cases = []SelectTestCase{
		{0, 2},
		{1, 3},
		{60, 62},
		{61, 66},
		{121, 126},
		{122, 129},
		{183, 190},
		{184, 193},
		{245, 254},
		{246, 257},
		{307, 318},
		{308, 321},
		{369, 382},
		{370, 385},
		{431, 446},
		{432, 449},
		{493, 510},
		{494, 513},
		{1005, 1024},
		{5980, 5998},
	}
)

func TestHasSelectIndex(t *testing.T) {
	builder := NewVectorBuilder()

	for _, v := range bitCases {
		builder.Set(v.pos, v.bit)
	}

	for _, v := range bitCases {
		x, err := builder.Get(v.pos)
		if err != nil || x != v.bit {
			t.Error("Expected", v.bit, "got", x)
		}
	}

	for _, v := range bitCases {
		x, err := builder.Get(v.pos)
		if err != nil || x != v.bit {
			t.Error("Expected", v.bit, "got", x)
		}
	}
	
	size := builder.Size()
	if size != uint64(6001) {
		t.Error("Expected", 6001, "got", size)
	}

	vec, _ := builder.Build(true, true)

	for _, v := range bitCases {
		x, err := vec.Get(v.pos)
		if err != nil || x != v.bit {
			t.Error("Expected", v.bit, "got", x)
		}
	}

	size = vec.Size()
	if size != uint64(6001) {
		t.Error("Expected", 6001, "got", size)
	}

	size = vec.NumOfBits(true)
	if size != uint64(20) {
		t.Error("Expected", 20, "got", size)
	}

	size = vec.NumOfBits(false)
	if size != uint64(5981) {
		t.Error("Expected", 5981, "got", size)
	}

	for _, v := range rankCases {
		rank, err := vec.Rank1(v.pos)
		if err != nil || rank != v.rank {
			t.Error("Expected", v.rank, "got", rank)
		}
	}

	for _, v := range rankCases {
		rank, err := vec.Rank0(v.pos)
		if err != nil || rank != (v.pos-v.rank) {
			t.Error("Expected", (v.pos - v.rank), "got", rank)
		}
	}

	for _, v := range select1Cases {
		pos, err := vec.Select1(v.index)
		if err != nil || pos != v.pos {
			t.Error("Expected", v.pos, "got", pos)
		}
	}

	for _, v := range select0Cases {
		pos, err := vec.Select0(v.index)
		if err != nil || pos != v.pos {
			t.Error("Expected", v.pos, "got", pos)
		}
	}

}

func TestNoSelectIndex(t *testing.T) {
	builder := NewVectorBuilder()

	for _, v := range bitCases {
		builder.Set(v.pos, v.bit)
	}

	vec, _ := builder.Build(false, false)

	for _, v := range bitCases {
		x, err := vec.Get(v.pos)
		if err != nil || x != v.bit {
			t.Error("Expected", v.bit, "got", x)
		}
	}

	for _, v := range rankCases {
		rank, err := vec.Rank1(v.pos)
		if err != nil || rank != v.rank {
			t.Error("Expected", v.rank, "got", rank)
		}
	}

	for _, v := range select1Cases {
		pos, err := vec.Select1(v.index)
		if err != nil || pos != v.pos {
			t.Error("Expected", v.pos, "got", pos)
		}
	}

	for _, v := range select0Cases {
		pos, err := vec.Select0(v.index)
		if err != nil || pos != v.pos {
			t.Error("Expected", v.pos, "got", pos)
		}
	}
}

func TestOutOfRange(t *testing.T) {
	builder := NewVectorBuilder()

	for _, v := range bitCases {
		builder.Set(v.pos, v.bit)
	}

	vec, err := builder.Build(true, true)

	bit, err := vec.Get(6002)
	if err == nil || err != ErrorOutOfRange || bit == true {
		t.Error()
	}

	rank, err := vec.Rank(6002, true)
	if err == nil || err != ErrorOutOfRange || rank != NotFound {
		t.Error()
	}

	rank, err = vec.Rank(6002, false)
	if err == nil || err != ErrorOutOfRange || rank != NotFound {
		t.Error()
	}

	pos, err := vec.Select(20, true)
	if err == nil || err != ErrorOutOfRange || pos != NotFound {
		t.Error()
	}

	pos, err = vec.Select(5981, false)
	if err == nil || err != ErrorOutOfRange || pos != NotFound {
		t.Error()
	}
}

func TestPushBack(t *testing.T) {
	builder := NewVectorBuilder()

	builder.PushBack(true)
	builder.PushBack(false)
	builder.PushBack(true)
	builder.PushBack(true)

	vec, err := builder.Build(false, false)

	pos, err := vec.Select1(2)
	if err != nil || pos != 3 {
		t.Error("Expected", 3, "got", pos)
	}
}

func TestMultiBits(t *testing.T) {
	builder := NewVectorBuilder()
	builder.PushBackBits(0x00FFFFFFFFFFFFFF, 63)
	builder.PushBackBits(0xFF55, 8)

	x, err := builder.GetBits(71, 1)
	if err == nil || err != ErrorOutOfRange {
		t.Error()
	}

	x, err = builder.GetBits(0, 64)
	if err != nil || x != uint64(0x80ffffffffffffff) {
		t.Error("Expected", uint64(0x80FFFFFFFFFFFFFF), "got", x)
	}

	x, err = builder.GetBits(8, 63)
	if err != nil || x != uint64(0x2A80FFFFFFFFFFFF) {
		t.Errorf("%x", x)
	}

	vec, err := builder.Build(true, true)

	size := vec.Size()
	if size != 71 {
		t.Error()
	}

	x, err = vec.GetBits(71, 1)
	if err == nil || err != ErrorOutOfRange {
		t.Error()
	}

	x, err = vec.GetBits(0, 64)
	if err != nil || x != uint64(0x80ffffffffffffff) {
		t.Error("Expected", uint64(0x80FFFFFFFFFFFFFF), "got", x)
	}

	x, err = vec.GetBits(8, 63)
	if err != nil || x != uint64(0x2A80FFFFFFFFFFFF) {
		t.Errorf("%x", x)
	}

	pos, err := vec.Select1(60)
	if err != nil || pos != 71 {
		t.Error()
	}
}

func TestDenseVector(t *testing.T) {
	builder := NewVectorBuilder()
	for i := uint64(0); i < uint64(0xFFF); i++ {
		builder.PushBack(true)
	}
	vec, err := builder.Build(false, false)
	pos, err := vec.Select1(513)
	if err != nil || pos != 513 {
		t.Error(pos)
	}
}

func TestMarshal(t *testing.T) {
	builder := NewVectorBuilder()

	for _, v := range bitCases {
		builder.Set(v.pos, v.bit)
	}

	vec, err := builder.Build(true, true)

	buffer, err := vec.MarshalBinary()
	if err != nil || len(buffer) == 0 {
		t.Error()
	}

	vec2, err := NewVectorFromBinary(buffer)
	if err != nil {
		t.Error(err)
	}
	for _, v := range bitCases {
		x, err := vec2.Get(v.pos)
		if err != nil || x != v.bit {
			t.Error("Expected", v.bit, "got", x)
		}
	}

	size := vec2.Size()
	if size != uint64(6001) {
		t.Error("Expected", 6001, "got", size)
	}

	size = vec2.NumOfBits(true)
	if size != uint64(20) {
		t.Error("Expected", 20, "got", size)
	}

	size = vec2.NumOfBits(false)
	if size != uint64(5981) {
		t.Error("Expected", 5981, "got", size)
	}

	for _, v := range rankCases {
		rank, err := vec2.Rank1(v.pos)
		if err != nil || rank != v.rank {
			t.Error("Expected", v.rank, "got", rank)
		}
	}

	for _, v := range rankCases {
		rank, err := vec2.Rank0(v.pos)
		if err != nil || rank != (v.pos-v.rank) {
			t.Error("Expected", (v.pos - v.rank), "got", rank)
		}
	}

	for _, v := range select1Cases {
		pos, err := vec2.Select1(v.index)
		if err != nil || pos != v.pos {
			t.Error("Expected", v.pos, "got", pos)
		}
	}

	for _, v := range select0Cases {
		pos, err := vec2.Select0(v.index)
		if err != nil || pos != v.pos {
			t.Error("Expected", v.pos, "got", pos)
		}
	}

	var buf []byte
	vec3, err := NewVectorFromBinary(buf)
	if err == nil || err != ErrorInvalidLength {
		t.Error(err.Error())
	}
	buf = make([]byte, minimumSize+1)
	binary.LittleEndian.PutUint64(buf, uint64(minimumSize))
	err = vec3.UnmarshalBinary(buf)
	if err == nil || err != ErrorInvalidLength {
		t.Error(err.Error())
	}

	var ri rankIndex
	ribuf := make([]byte, 15)
	err = ri.UnmarshalBinary(ribuf)
	if err == nil {
		t.Error(err.Error())
	}

	builder = NewVectorBuilderWithInit(vec3)
	builder.PushBack(true)
	vec3, err = builder.Build(true, true)
	buffer, err = vec3.MarshalBinary()
	badBuf := make([]byte, len(buffer))
	copy(badBuf, buffer)
	badBuf[24] = 0xFF
	err = vec3.UnmarshalBinary(badBuf)
	if err != ErrorInvalidFormat {
		t.Error()
	}

	copy(badBuf, buffer)
	badBuf[36] = 0xFF
	err = vec3.UnmarshalBinary(badBuf)
	if err != ErrorInvalidFormat {
		t.Error()
	}

	copy(badBuf, buffer)
	badBuf[72] = 0xFF
	err = vec3.UnmarshalBinary(badBuf)
	if err != ErrorInvalidFormat {
		t.Error()
	}
	copy(badBuf, buffer)
	badBuf[92] = 0xFF
	err = vec3.UnmarshalBinary(badBuf)
	if err != ErrorInvalidFormat {
		t.Error()
	}
}
