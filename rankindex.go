package sbvector

const (
	mask7F uint64 = 0x7F
	maskFF uint64 = 0xFF
	mask1F uint64 = 0x1FF
)

// RankIndex holds number of non-zero bits at each ranges.
type RankIndex struct {
	absVal uint64
	rel    uint64
}

type rankIndexMethods interface {
	abs() uint64
	rel1() uint64
	rel2() uint64
	rel3() uint64
	rel4() uint64
	rel5() uint64
	rel6() uint64
	rel7() uint64

	setAbs(val uint64)
	setRel1(val uint64)
	setRel2(val uint64)
	setRel3(val uint64)
	setRel4(val uint64)
	setRel5(val uint64)
	setRel6(val uint64)
	setRel7(val uint64)
}

func (index *RankIndex) abs() uint64 {
	return index.absVal
}

func (index *RankIndex) rel1() uint64 {
	return (index.rel & mask7F)
}

func (index *RankIndex) rel2() uint64 {
	return ((index.rel >> 7) & maskFF)
}

func (index *RankIndex) rel3() uint64 {
	return ((index.rel >> 15) & maskFF)
}

func (index *RankIndex) rel4() uint64 {
	return ((index.rel >> 23) & mask1F)
}

func (index *RankIndex) rel5() uint64 {
	return ((index.rel >> 32) & mask1F)
}

func (index *RankIndex) rel6() uint64 {
	return ((index.rel >> 41) & mask1F)
}

func (index *RankIndex) rel7() uint64 {
	return ((index.rel >> 50) & mask1F)
}

func (index *RankIndex) setAbs(val uint64) {
	index.absVal = val
}

func (index *RankIndex) setRel1(val uint64) {
	index.rel = ((index.rel & ^mask7F) | (val & mask7F))
}

func (index *RankIndex) setRel2(val uint64) {
	index.rel = ((index.rel & ^(maskFF << 7)) | ((val & maskFF) << 7))
}

func (index *RankIndex) setRel3(val uint64) {
	index.rel = ((index.rel & ^(maskFF << 15)) | ((val & maskFF) << 15))
}

func (index *RankIndex) setRel4(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 23)) | ((val & mask1F) << 23))
}

func (index *RankIndex) setRel5(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 32)) | ((val & mask1F) << 32))
}

func (index *RankIndex) setRel6(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 41)) | ((val & mask1F) << 41))
}

func (index *RankIndex) setRel7(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 50)) | ((val & mask1F) << 50))
}
