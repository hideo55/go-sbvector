package sbvector

const (
	mask7F uint64 = 0x7F
	maskFF uint64 = 0xFF
	mask1F uint64 = 0x1FF
)

type rankindex struct {
	absVal uint64
	rel uint64
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

func (index *rankindex) abs() uint64 {
	return index.absVal
}

func (index *rankindex) rel1() uint64 {
	return (index.rel & mask7F)
}

func (index *rankindex) rel2() uint64 {
	return ((index.rel >> 7) & maskFF)
}

func (index *rankindex) rel3() uint64 {
	return ((index.rel >> 15) & maskFF)
}

func (index *rankindex) rel4() uint64 {
	return ((index.rel >> 23) & mask1F)
}

func (index *rankindex) rel5() uint64 {
	return ((index.rel >> 32) & mask1F)
}

func (index *rankindex) rel6() uint64 {
	return ((index.rel >> 41) & mask1F)
}

func (index *rankindex) rel7() uint64 {
	return ((index.rel >> 50) & mask1F)
}

func (index *rankindex) setAbs(val uint64) {
	index.absVal = val
}

func (index *rankindex) setRel1(val uint64) {
	index.rel = ((index.rel & ^mask7F) | (val & mask7F))
}

func (index *rankindex) setRel2(val uint64) {
	index.rel = ((index.rel & ^(maskFF << 7)) | ((val & maskFF) << 7))
}

func (index *rankindex) setRel3(val uint64) {
	index.rel = ((index.rel & ^(maskFF << 15)) | ((val & maskFF) << 15))
}

func (index *rankindex) setRel4(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 23)) | ((val & mask1F) << 23))
}

func (index *rankindex) setRel5(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 32)) | ((val & mask1F) << 32))
}

func (index *rankindex) setRel6(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 41)) | ((val & mask1F) << 41))
}

func (index *rankindex) setRel7(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 50)) | ((val & mask1F) << 50))
}
