package secp256k1

import (
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
		return a.Equal(a)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestIdentityIsIdentity(t *testing.T) {
	if !NewPoint().IsIdentity() {
		t.Error("NewPoint() didn't return identity point")
	}
}

func TestPointAdditionCommutative(t *testing.T) {
	err := quick.Check(func(a *Point, b *Point) bool {
		way1 := a.Add(b)
		way2 := b.Add(a)
		return way1.Equal(way2)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestPointAddIdentityDoesNothing(t *testing.T) {
	err := quick.Check(func(a *Point) bool {
		return a.Add(NewPoint()).Equal(a)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestPointSelfSubtractionIsIdentity(t *testing.T) {
	err := quick.Check(func(a *Point) bool {
		return a.Sub(a).IsIdentity()
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestPointSubtractionIsAddNegated(t *testing.T) {
	err := quick.Check(func(a *Point, b *Point) bool {
		way1 := a.Sub(b)
		way2 := a.Add(b.Negate())
		return way1.Equal(way2)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}
