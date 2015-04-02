package sbvector

// BitVectorBuilderData holds bit vector data to build.
type BitVectorBuilderData struct {
	vec *bitVectorData
}

// SuccinctBitVectorBuilder is interface of succinct bit vector builder.
type SuccinctBitVectorBuilder interface {
	Set(i uint64, val bool)
	Get(i uint64) (bool, error)
	PushBack(b bool)
	PushBackBits(x uint64, length uint64)
	GetBits(pos uint64, length uint64) (uint64, error)
	Size() uint64
	Build(enableFasterSelect1 bool, enableFasterSelect0 bool) (SuccinctBitVector, error)
}

// NewVectorBuilder returns new succinct bit vector builder.
func NewVectorBuilder() SuccinctBitVectorBuilder {
	builder := new(BitVectorBuilderData)
	builder.vec = new(bitVectorData)
	return builder
}

// NewVectorBuilderWithInit returns new succinct bit vector builder(initialize by argument).
func NewVectorBuilderWithInit(vec *bitVectorData) SuccinctBitVectorBuilder {
	builder := new(BitVectorBuilderData)
	builder.vec = vec
	return builder
}

// Set value to bit vector by index
func (builder *BitVectorBuilderData) Set(i uint64, val bool) {
	builder.vec.set(i, val)
}

// Get returns value from bit vector by index.
func (builder *BitVectorBuilderData) Get(i uint64) (bool, error) {
	return builder.vec.Get(i)
}

// PushBack add bit to the bit vector
func (builder *BitVectorBuilderData) PushBack(b bool) {
	builder.vec.pushBack(b)
}

// PushBackBits add bits to the bit vector
func (builder *BitVectorBuilderData) PushBackBits(x uint64, length uint64) {
	builder.vec.pushBackBits(x, length)
}

//GetBits returns bits from bit vector
func (builder *BitVectorBuilderData) GetBits(pos uint64, length uint64) (uint64, error) {
	return builder.vec.GetBits(pos, length)
}

// Size returns size of bit vector
func (builder *BitVectorBuilderData) Size() uint64 {
	return builder.vec.Size()
}

// Build creates indexes for succinct bit vector(rank index, ...).
// If `enableFasterSelect1` is true, creates index for select1 make faster.
// If `enableFasterSelect0` is true, creates index for select0 make faster.
func (builder *BitVectorBuilderData) Build(enableFasterSelect1 bool, enableFasterSelect0 bool) (SuccinctBitVector, error) {
	builder.vec.build(enableFasterSelect1, enableFasterSelect0)
	vec := builder.vec
	builder.vec = new(bitVectorData)
	return vec, nil
}
