package sbvector

import (
	"testing"
)

func TestSBVector(t *testing.T) {
    a := SBVector{}
    a.set(65, true)
	a.set(200, true)
	var x = a.get(65)
    if x != true {
		t.Error("Expected", true, "got", x)
	}
	a.build(false, false)
	var rank = a.rank(66)
	if rank != 1 {
		t.Error()
	}
	if a.select1(0) != 65 {
		t.Error()
	}
	if a.select1(1) != 200 {
		t.Error()
	}
}

