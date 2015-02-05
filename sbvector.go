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
		{255, true},
		{256, true},
		{300, false},
	}

	rankCases = []RankTestCase{
		{1, 1},
		{2, 2},
		{3, 2},
		{63, 2},
		{64, 3},
		{65, 4},
		{66, 5},
		{255, 5},
		{256, 6},
		{257, 7},
		{300, 7},
		{301, 7},
	}

	select1Cases = []SelectTestCase{
		{0, 0},
		{1, 1},
		{2, 63},
		{3, 64},
		{4, 65},
		{5, 255},
		{6, 256},
	}

	select0Cases = []SelectTestCase{
		{0, 2},
		{1, 3},
		{60, 62},
		{61, 66},
		{249, 254},
		{250, 257},
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
