package secp256k1

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func (*Point) Generate(r *rand.Rand, size int) reflect.Value {
	// The idea is to generate x coordinates until we end up on the curve,
	// or give up and return the identity point.

	// Most of the coordinates we generate will end up on the curve, so a low value is sensible
	const generateIterations = 3

	for i := 0; i < generateIterations; i++ {
		x := randomFieldElement(r, size)
		fx := NewField().Set(x).Square().Mul(x).AddU64(b)
		if fx.HasSqrt() != 1 {
			continue
		}
		y := fx.Sqrt()
		z := NewField().SetUint64(1)
		return reflect.ValueOf(&Point{x, y, z, true})
	}
	return reflect.ValueOf(NewPoint())
}

func TestPointEqualToItself(t *testing.T) {
	err := quick.Check(func(a *Point) bool {
		fmt.Println("a", a)
		return a.Equal(a)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}
