package sbvector
 
const (
    MASK_7F uint64 = 0x7F
    MASK_FF uint64 = 0xFF
    MASK_1F uint64 = 0x1FF
)
 
type rankindex struct {
    abs_ uint64
    rel_ uint64
}
 
type  rankIndexMethods interface {
    abs() uint64
    rel1() uint64
    rel2() uint64
    rel3() uint64
    rel4() uint64
    rel5() uint64
    rel6() uint64
    rel7() uint64
 
    set_abs(val uint64)
    set_rel1(val uint64)
    set_rel2(val uint64)
    set_rel3(val uint64)
    set_rel4(val uint64)
    set_rel5(val uint64)
    set_rel6(val uint64)
    set_rel7(val uint64)
}
 
func (this *rankindex) abs() uint64 {
    return this.abs_
}
 
func (this *rankindex) rel1() uint64 {
    return ( this.rel_ & MASK_7F )
}
 
func (this *rankindex) rel2() uint64 {
    return ( ( this.rel_ >> 7 ) & MASK_FF )
}
 
func (this *rankindex) rel3() uint64 {
    return ( ( this.rel_ >> 16 ) & MASK_FF )
}
 
func (this *rankindex) rel4() uint64 {
    return ( ( this.rel_ >> 23 ) & MASK_1F )
}
 
func (this *rankindex) rel5() uint64 {
    return ( ( this.rel_ >> 32 ) & MASK_1F )
}
 
func (this *rankindex) rel6() uint64 {
    return ( ( this.rel_ >> 41 ) & MASK_1F )
}
 
func (this *rankindex) rel7() uint64 {
    return ( ( this.rel_ >> 50 ) & MASK_1F )
}
 
func (this *rankindex) set_abs(val uint64) {
    this.abs_ = val
}
 
func (this *rankindex) set_rel1(val uint64) {
    this.rel_ = ((this.rel_ & ^MASK_7F) | (val & MASK_7F))
}
 
func (this *rankindex) set_rel2(val uint64) {
    this.rel_ = ((this.rel_ & ^(MASK_FF << 7)) | ((val & MASK_FF) << 7))
}
 
func (this *rankindex) set_rel3(val uint64) {
    this.rel_ = ((this.rel_ & ^(MASK_FF << 15)) | ((val & MASK_FF) << 15))
}
 
func (this *rankindex) set_rel4(val uint64) {
    this.rel_ = ((this.rel_ & ^(MASK_1F << 23)) | ((val & MASK_1F) << 23))
}
 
func (this *rankindex) set_rel5(val uint64) {
    this.rel_ = ((this.rel_ & ^(MASK_1F << 32)) | ((val & MASK_1F) << 32))
}
 
func (this *rankindex) set_rel6(val uint64) {
    this.rel_ = ((this.rel_ & ^(MASK_1F << 41)) | ((val & MASK_1F) << 41))
}
 
func (this *rankindex) set_rel7(val uint64) {
    this.rel_ = ((this.rel_ & ^(MASK_1F << 50)) | ((val & MASK_1F) << 50))
}

