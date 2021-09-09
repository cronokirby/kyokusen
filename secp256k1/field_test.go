package secp256k1

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func (*Field) Generate(r *rand.Rand, size int) reflect.Value {
	out := NewField()
	data := make([]byte, FieldBytes)
	// Fill in a certain number of bytes with zero. Smaller sizes will be closer to zero.
	for i := 0; i < size && i < len(data); i++ {
		data[len(data)-i-1] = byte(r.Uint32())
	}
	_ = out.UnmarshalBinary(data)
	return reflect.ValueOf(out)
}

func TestFieldAdditionCommutative(t *testing.T) {
	err := quick.Check(func(a *Field, b *Field) bool {
		way1 := NewField().Set(a).Add(b)
		way2 := NewField().Set(b).Add(a)
		return way1.Eq(way2) == 1
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestFieldAddZeroIdentity(t *testing.T) {
	err := quick.Check(func(a *Field) bool {
		shouldBeA := NewField().Add(a)
		return shouldBeA.Eq(a) == 1
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestFieldMultiplicationCommutative(t *testing.T) {
	err := quick.Check(func(a *Field, b *Field) bool {
		way1 := NewField().Set(a).Mul(b)
		way2 := NewField().Set(b).Mul(a)
		return way1.Eq(way2) == 1
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestMultiplyOneIdentity(t *testing.T) {
	err := quick.Check(func(a *Field) bool {
		shouldBeA := NewField().SetUint64(1).Mul(a)
		return shouldBeA.Eq(a) == 1
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestMultiplyInverse(t *testing.T) {
	err := quick.Check(func(a *Field) bool {
		shouldBeOne := NewField().Set(a).Invert().Mul(a)
		one := NewField().SetUint64(1)
		return shouldBeOne.Eq(one) == 1
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}
