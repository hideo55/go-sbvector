package sbvector

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
)

const (
	mask7F     uint64 = 0x7F
	maskFF     uint64 = 0xFF
	mask1F     uint64 = 0x1FF
	binarySize int    = 16
)

// rankIndex holds number of non-zero bits at each ranges.
type rankIndex struct {
	absVal uint64
	rel    uint64
}

type rankIndexMethods interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

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

func (index *rankIndex) abs() uint64 {
	return index.absVal
}

func (index *rankIndex) rel1() uint64 {
	return (index.rel & mask7F)
}

func (index *rankIndex) rel2() uint64 {
	return ((index.rel >> 7) & maskFF)
}

func (index *rankIndex) rel3() uint64 {
	return ((index.rel >> 15) & maskFF)
}

func (index *rankIndex) rel4() uint64 {
	return ((index.rel >> 23) & mask1F)
}

func (index *rankIndex) rel5() uint64 {
	return ((index.rel >> 32) & mask1F)
}

func (index *rankIndex) rel6() uint64 {
	return ((index.rel >> 41) & mask1F)
}

func (index *rankIndex) rel7() uint64 {
	return ((index.rel >> 50) & mask1F)
}

func (index *rankIndex) setAbs(val uint64) {
	index.absVal = val
}

func (index *rankIndex) setRel1(val uint64) {
	index.rel = ((index.rel & ^mask7F) | (val & mask7F))
}

func (index *rankIndex) setRel2(val uint64) {
	index.rel = ((index.rel & ^(maskFF << 7)) | ((val & maskFF) << 7))
}

func (index *rankIndex) setRel3(val uint64) {
	index.rel = ((index.rel & ^(maskFF << 15)) | ((val & maskFF) << 15))
}

func (index *rankIndex) setRel4(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 23)) | ((val & mask1F) << 23))
}

func (index *rankIndex) setRel5(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 32)) | ((val & mask1F) << 32))
}

func (index *rankIndex) setRel6(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 41)) | ((val & mask1F) << 41))
}

func (index *rankIndex) setRel7(val uint64) {
	index.rel = ((index.rel & ^(mask1F << 50)) | ((val & mask1F) << 50))
}

func (index *rankIndex) MarshalBinary() ([]byte, error) {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, &index.absVal)
	binary.Write(buffer, binary.LittleEndian, &index.rel)
	return buffer.Bytes(), nil
}

func (index *rankIndex) UnmarshalBinary(data []byte) error {
	buf := data
	if len(buf) != binarySize {
		return errors.New("Invalid value")
	}
	buf = data[:8]
	index.absVal = binary.LittleEndian.Uint64(buf)
	buf = data[8:]
	index.rel = binary.LittleEndian.Uint64(buf)
	return nil
}
