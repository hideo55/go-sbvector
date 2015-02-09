package sbvector

import (
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
	}
)

func TestHasSelectIndex(t *testing.T) {
	vec, err := NewVector()
	if err != nil {
		t.Error()
	}

	for _, v := range bitCases {
		vec.Set(v.pos, v.bit)
	}

	for _, v := range bitCases {
		x, err := vec.Get(v.pos)
		if err != nil || x != v.bit {
			t.Error("Expected", v.bit, "got", x)
		}
	}

	vec.Build(true, true)

	for _, v := range bitCases {
		x, err := vec.Get(v.pos)
		if err != nil || x != v.bit {
			t.Error("Expected", v.bit, "got", x)
		}
	}

	size := vec.Size()
	if size != uint64(1025) {
		t.Error("Expected", 1025, "got", size)
	}

	size = vec.NumOfBits(true)
	if size != uint64(19) {
		t.Error("Expected", 19, "got", size)
	}

	size = vec.NumOfBits(false)
	if size != uint64(1006) {
		t.Error("Expected", 106, "got", size)
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

func TestNoSelectIndex(t *testing.T) {
	vec, err := NewVector()
	if err != nil {
		t.Error()
	}

	for _, v := range bitCases {
		vec.Set(v.pos, v.bit)
	}

	vec.Build(false, false)

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
	vec, err := NewVector()
	if err != nil {
		t.Error()
	}

	for _, v := range bitCases {
		vec.Set(v.pos, v.bit)
	}

	vec.Build(true, true)

	bit, err := vec.Get(1026)
	if err == nil || bit == true {
		t.Error()
	}

	rank, err := vec.Rank(1026, true)
	if err == nil || rank != NotFound {
		t.Error()
	}

	rank, err = vec.Rank(1026, false)
	if err == nil || rank != NotFound {
		t.Error()
	}

	pos, err := vec.Select(19, true)
	if err == nil || pos != NotFound {
		t.Error()
	}

	pos, err = vec.Select(1006, false)
	if err == nil || pos != NotFound {
		t.Error()
	}
}

func TestPushBack(t *testing.T) {
	vec, err := NewVector()
	if err != nil {
		t.Error()
	}

	vec.PushBack(true)
	vec.PushBack(false)
	vec.PushBack(true)
	vec.PushBack(true)

	vec.Build(false, false)

	pos, err := vec.Select1(2)
	if err != nil || pos != 3 {
		t.Error("Expected", 3, "got", pos)
	}
}

func TestMultiBits(t *testing.T) {
	vec, err := NewVector()
	if err != nil {
		t.Error()
	}
	vec.PushBackBits(0x00FFFFFFFFFFFFFF, 63)
	vec.PushBackBits(0xFF55, 8)
	vec.Build(true, true)

	size := vec.Size()
	if size != 71 {
		t.Error()
	}

	x, err := vec.GetBits(71, 1)
	if err == nil {
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
