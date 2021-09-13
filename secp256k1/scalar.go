package secp256k1

import (
	"errors"

	"github.com/cronokirby/kyokusen"
	"github.com/cronokirby/saferith"
)

// q is the modulus for the scalars used in secp256k1.
var q *saferith.Modulus

func init() {
	q, _ = saferith.ModulusFromHex("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141")
}

type Scalar struct {
	nat saferith.Nat
}

// MarshalBinary returns the contents of this scalar as Big Endian bytes.
func (s *Scalar) MarshalBinary() ([]byte, error) {
	// This makes sure that the size of the scalar is padded to the right size.
	s.nat.Mod(&s.nat, q)
	return s.nat.Bytes(), nil
}

// UnmarshalBinary deserializes Big Endian bytes into this scalar.
func (s *Scalar) UnmarshalBinary(data []byte) error {
	if len(data) != p.BitLen()/8 {
		return errors.New("secp256k1.Scalar.UnmarshalBinary: invalid data length")
	}
	s.nat.SetBytes(data)
	if _, _, lt := s.nat.CmpMod(p); lt != 1 {
		return errors.New("secp256k1.Scalar.UnmarshalBinary: value is greater than order")
	}
	return nil
}

// Curve returns the curve associated with this scalar field.
func (s *Scalar) Curve() kyokusen.Curve {
	return nil
}

func (s1 *Scalar) Add(other kyokusen.Scalar) kyokusen.Scalar {
	return nil
}

func (s1 *Scalar) Sub(other kyokusen.Scalar) kyokusen.Scalar {
	return nil
}

func (s1 *Scalar) Negate() kyokusen.Scalar {
	return nil
}

func (s1 *Scalar) Mul(other kyokusen.Scalar) kyokusen.Scalar {
	return nil
}

func (s1 *Scalar) Invert() kyokusen.Scalar {
	return nil
}

func (s1 *Scalar) Equal(other kyokusen.Scalar) bool {
	return false
}

func (s1 *Scalar) IsZero() bool {
	return false
}

func (s1 *Scalar) Set(other kyokusen.Scalar) kyokusen.Scalar {
	return nil
}

func (s1 *Scalar) SetNat(other *saferith.Nat) kyokusen.Scalar {
	return nil
}

func (s1 *Scalar) Act(kyokusen.Point) kyokusen.Point {
	return nil
}

func (s1 *Scalar) ActOnBase() kyokusen.Point {
	return nil
}
