package secp256k1

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/cronokirby/saferith"
)

func randomScalar(r *rand.Rand, size int) *Scalar {
	data := make([]byte, q.BitLen()/8)
	// Fill in a certain number of bytes with zero. Smaller sizes will be closer to zero.
	for i := 0; i < size && i < len(data); i++ {
		data[len(data)-i-1] = byte(r.Uint32())
	}
	return NewScalar().SetNat(new(saferith.Nat).SetBytes(data)).(*Scalar)
}

func (*Scalar) Generate(r *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(randomScalar(r, size))
}

func TestScalarAdditionCommutative(t *testing.T) {
	err := quick.Check(func(a, b *Scalar) bool {
		way1 := NewScalar().Set(a).Add(b)
		way2 := NewScalar().Set(b).Add(a)
		return way1.Equal(way2)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestScalarAddZeroIdentity(t *testing.T) {
	err := quick.Check(func(a *Scalar) bool {
		shouldBeA := NewScalar().Add(a)
		return shouldBeA.Equal(a)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestScalarMultiplicationCommutative(t *testing.T) {
	err := quick.Check(func(a, b *Scalar) bool {
		way1 := NewScalar().Set(a).Mul(b)
		way2 := NewScalar().Set(b).Mul(a)
		return way1.Equal(way2)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestScalarMultiplyOneIdentity(t *testing.T) {
	err := quick.Check(func(a *Scalar) bool {
		shouldBeA := NewScalar().SetNat(new(saferith.Nat).SetUint64(1)).Mul(a)
		return shouldBeA.Equal(a)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestScalarMultiplyInverse(t *testing.T) {
	err := quick.Check(func(a *Scalar) bool {
		shouldBeOne := NewScalar().Set(a).Invert().Mul(a)
		one := NewScalar().SetNat(new(saferith.Nat).SetUint64(1))
		return shouldBeOne.Equal(one)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}

func TestScalarMarshalRoundtrip(t *testing.T) {
	err := quick.Check(func(a *Scalar) bool {
		marshalled, err := a.MarshalBinary()
		if err != nil {
			return false
		}
		unmarshalled := NewScalar()
		if err := unmarshalled.UnmarshalBinary(marshalled); err != nil {
			return false
		}
		return unmarshalled.Equal(a)
	}, &quick.Config{})
	if err != nil {
		t.Error(err)
	}
}
