package sbvector
 
import "testing"
import "fmt"
 
func TestSBVector(t *testing.T) {
    a := sbvector{}
    a.set(65, true)
    fmt.Printf("%d\n", a.get(65))
}

