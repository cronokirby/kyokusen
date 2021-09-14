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

func (s *Scalar) String() string {
	return s.nat.String()
}

func NewScalar() *Scalar {
	return &Scalar{}
}

// castScalar converts a scalar implementing the generic interface to this specific type.
//
// Since implementors of the Scalar interface are only expected to work with
// their own type, we are allowed to cast the interface at the beginning of our methods.
func castScalar(s kyokusen.Scalar) *Scalar {
	casted, ok := s.(*Scalar)
	if !ok {
		panic("failed to cast type to *secp256k1.Scalar")
	}
	return casted
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
	s2 := castScalar(other)
	s1.nat.ModAdd(&s1.nat, &s2.nat, q)
	return s1
}

func (s1 *Scalar) Sub(other kyokusen.Scalar) kyokusen.Scalar {
	s2 := castScalar(other)
	s1.nat.ModSub(&s1.nat, &s2.nat, q)
	return s1
}

func (s1 *Scalar) Negate() kyokusen.Scalar {
	s1.nat.ModNeg(&s1.nat, q)
	return s1
}

func (s1 *Scalar) Mul(other kyokusen.Scalar) kyokusen.Scalar {
	s2 := castScalar(other)
	s1.nat.ModMul(&s1.nat, &s2.nat, q)
	return s1
}

func (s1 *Scalar) Invert() kyokusen.Scalar {
	s1.nat.ModInverse(&s1.nat, q)
	return s1
}

func (s1 *Scalar) Equal(other kyokusen.Scalar) bool {
	s2 := castScalar(other)
	return s1.nat.Eq(&s2.nat) == 1
}

func (s1 *Scalar) IsZero() bool {
	return s1.nat.EqZero() == 1
}

func (s1 *Scalar) Set(other kyokusen.Scalar) kyokusen.Scalar {
	s2 := castScalar(other)
	s1.nat.SetNat(&s2.nat)
	return s1
}

func (s1 *Scalar) SetNat(other *saferith.Nat) kyokusen.Scalar {
	s1.nat.Mod(other, q)
	return s1
}

func (s *Scalar) Act(other kyokusen.Point) kyokusen.Point {
	bytes := s.nat.Bytes()
	acc := NewPoint()
	for _, b := range bytes {
		for i := 7; i >= 0; i-- {
			acc = acc.Add(acc).(*Point)
			added := acc.Add(other).(*Point)
			acc.CondAssign(saferith.Choice((b>>i)&1), added)
		}
	}
	return acc
}

func (s *Scalar) ActOnBase() kyokusen.Point {
	return s.Act(basePoint)
}
