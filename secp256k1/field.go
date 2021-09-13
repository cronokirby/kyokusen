package secp256k1

import (
	"errors"

	"github.com/cronokirby/saferith"
)

// fieldBytes is the number of bytes in our encoding of a field element.
const fieldBytes = 32

// p is the modulus for the field used in secp256k1.
var p *saferith.Modulus

// pDiv2 is (p - 1) / 2, useful for checking if a value has a square root
var pDiv2 *saferith.Nat

func init() {
	pNat, _ := new(saferith.Nat).SetHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F")
	p = saferith.ModulusFromNat(pNat)
	pNat.Rsh(pNat, 1, p.BitLen())
	pDiv2 = pNat
}

// FieldBytes is the number of bytes in the field.
const FieldBytes = 32

// Field represents an element in the prime field used by secp256k1.
//
// This field is used later to implement point operations on the curve.
type Field struct {
	nat saferith.Nat
}

// NewField creates a new field element, with its value set to 0.
func NewField() *Field {
	var nat saferith.Nat
	// This will conveniently set the announced size, and mark this number as reduced modulo p.
	nat.Mod(&nat, p)
	return &Field{nat: nat}
}

// Set calculates z <- x, returning z.
func (z *Field) Set(x *Field) *Field {
	z.nat.SetNat(&x.nat)
	return z
}

// SetUint64 calculates z <- x, returning z.
func (z *Field) SetUint64(x uint64) *Field {
	z.nat.SetUint64(x)
	return z
}

// CondAssign sets z <- x, only if yes = 1, in constant-time.
func (z *Field) CondAssign(yes saferith.Choice, x *Field) *Field {
	z.nat.CondAssign(yes, &x.nat)
	return z
}

// CondNegate sets z <- -z, only if yes = 1, in constant-time.
func (z *Field) CondNegate(yes saferith.Choice) *Field {
	negated := NewField().Set(z).Negate()
	return z.CondAssign(yes, negated)
}

// String returns a string representation of this field element.
func (z *Field) String() string {
	return z.nat.String()
}

// Add calculates z <- z + a, returning z.
func (z *Field) Add(a *Field) *Field {
	z.nat.ModAdd(&z.nat, &a.nat, p)
	return z
}

// Add calculates z <- z + a, returning z.
//
// This may be faster than Add.
func (z *Field) AddU64(a uint64) *Field {
	z.nat.ModAdd(&z.nat, new(saferith.Nat).SetUint64(a), p)
	return z
}

// Sub calculates z <- z - a, returning z.
func (z *Field) Sub(a *Field) *Field {
	z.nat.ModSub(&z.nat, &a.nat, p)
	return z
}

// Sub calculates z <- -z, returning z.
func (z *Field) Negate() *Field {
	z.nat.ModNeg(&z.nat, p)
	return z
}

// Mul calculates z <- z * a, returning z.
func (z *Field) Mul(a *Field) *Field {
	z.nat.ModMul(&z.nat, &a.nat, p)
	return z
}

// MulU64 calculates z <- z * a, returning z.
//
// This is more efficient than Mul.
func (z *Field) MulU64(a uint64) *Field {
	z.nat.ModMul(&z.nat, new(saferith.Nat).SetUint64(a), p)
	return z
}

// Square calculates z <- z * z, returning z.
func (z *Field) Square() *Field {
	return z.Mul(z)
}

// Invert calculates z <- z^-1, returning z.
func (z *Field) Invert() *Field {
	z.nat.ModInverse(&z.nat, p)
	return z
}

// Eq checks if two field values are equal, in constant-time.
func (z *Field) Eq(x *Field) saferith.Choice {
	return z.nat.Eq(&x.nat)
}

// Eq checks if a field value is equal to 0, in constant-time.
func (z *Field) EqZero() saferith.Choice {
	return z.nat.EqZero()
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
//
// This encodes the field element as big endian bytes. The result will always occupy
// 32 bytes of space.
func (z *Field) MarshalBinary() ([]byte, error) {
	// We do this to make sure that the announced size is padded to the size of p,
	// in case z was created wihout calling NewField, or something like that.
	// Since z will always be reduced modulo p, this doesn't actually cost anything.
	z.nat.Mod(&z.nat, p)
	return z.nat.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
//
// This expects exactly 32 Big Endian bytes, and will also return an error if the
// resulting value is >= the field modulus.
func (z *Field) UnmarshalBinary(data []byte) error {
	if len(data) != p.BitLen()/8 {
		return errors.New("secp256k1.Field.UnmarshalBinary: invalid data length")
	}
	z.nat.SetBytes(data)
	if _, _, lt := z.nat.CmpMod(p); lt != 1 {
		return errors.New("secp256k1.Field.UnmarshalBinary: value is greater than field prime")
	}
	return nil
}

// HasSqrt checks if a field value has a valid square root.
func (z *Field) HasSqrt() saferith.Choice {
	check := new(saferith.Nat).Exp(&z.nat, pDiv2, p)
	one := new(saferith.Nat).SetUint64(1)
	return check.Eq(one) | check.EqZero()
}

// Sqrt calculates z <- sqrt(z), if such a value exists. Otherwise, the result is undefined.
func (z *Field) Sqrt() *Field {
	z.nat.ModSqrt(&z.nat, p)
	return z
}

// IsEven returns a choice indicating if a field element is even.
func (z *Field) IsEven() saferith.Choice {
	return saferith.Choice(z.nat.Byte(0) & 1)
}
