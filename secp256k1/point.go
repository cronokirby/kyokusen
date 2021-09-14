package secp256k1

import (
	"crypto/subtle"
	"errors"
	"fmt"

	"github.com/cronokirby/kyokusen"
	"github.com/cronokirby/saferith"
)

// The b constant for the elliptic curve
const b = 7

// Point represents a point on the secp256k1 curve.
type Point struct {
	// Internally, we represent this as a projective point (X : Y : Z).
	// This corresponds to the affine point (X / Z, Y / Z), except when Z = 0,
	// which corresponds to the point at infinity.
	x *Field
	y *Field
	z *Field
	// This is a flag indicating that this Point's value is normalized. This
	// should be set exclusively based on which methods are called on a point,
	// making it okay to branch on this value.
	normalized bool
}

// basePoint is the generator of our elliptic curve group.
var basePoint *Point

func init() {
	xNat, _ := new(saferith.Nat).SetHex("79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798")
	yNat, _ := new(saferith.Nat).SetHex("483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8")
	basePoint = &Point{
		x: &Field{nat: *xNat},
		y: &Field{nat: *yNat},
		z: &Field{nat: *new(saferith.Nat).SetUint64(1)},
	}
}

// castPoint converts a point implementing the generic interface to this specific type.
//
// Since implementors of the Point interface are only expected to work with
// their own type, we are allowed to cast the interface at the beginning of our methods.
func castPoint(p kyokusen.Point) *Point {
	casted, ok := p.(*Point)
	if !ok {
		panic("failed to cast type to *secp256k1.Point")
	}
	return casted
}

func (p *Point) normalize() {
	if p.normalized {
		return
	}
	// We start (X : Y : Z).
	// If Z != 0, then we want to get (X/Z : Y/Z : 1)
	// If Z == 1, then we want to get (0 : 1 : 0)
	zZero := p.z.EqZero()
	zInverse := NewField().Set(p.z).Invert()
	one := NewField().SetUint64(1)
	// If Z != 0, then this will set the right values for X and Y
	p.x.Mul(zInverse)
	p.y.Mul(zInverse)
	// If Z == 0, then X needs to be 0, and Y needs to be 1.
	p.x.CondAssign(zZero, p.z)
	p.y.CondAssign(zZero, one)
	// If Z != 0, then Z needs to become 1.
	p.z.CondAssign(1^zZero, one)
	// The result is now normalized.
	p.normalized = true
}

// NewPoint returns the secp256k1 identity point.
func NewPoint() *Point {
	// (0 : 1 : 0) is the point at infinity, in projective coordinates.
	return &Point{
		x: NewField(),
		y: NewField().SetUint64(1),
		z: NewField(),
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("[%v : %v : %v]", p.x, p.y, p.z)
}

// MarshalBinary marshals a Secp256k1 point in the same way as Bitcoin does.
//
// The point at infinity can't be marshalled.
func (p *Point) MarshalBinary() ([]byte, error) {
	p.normalize()
	if p.IsIdentity() {
		return nil, errors.New("secp256k1: can't marshal point at infinity")
	}
	xBytes, err := p.x.MarshalBinary()
	if err != nil {
		return nil, err
	}
	out := make([]byte, 0, 1+len(xBytes))
	out = append(out, 3-byte(p.y.IsEven()))
	out = append(out, xBytes...)
	return out, nil
}

// UnmarshalBinary unmarshals a Secp256k1 point from Bitcoin's encoding.
func (p *Point) UnmarshalBinary(data []byte) error {
	if len(data) != 1+fieldBytes {
		return errors.New("secp256k1.UnmarshalBinary: invalid data")
	}
	if err := p.x.UnmarshalBinary(data[1:]); err != nil {
		return err
	}
	p.y.Set(p.x).Square().Mul(p.x).AddU64(b)
	if p.y.HasSqrt() != 1 {
		return errors.New("secp256k1.UnmarshalBinary: invalid point")
	}
	p.y.Sqrt()
	yShouldBeEven := saferith.Choice(subtle.ConstantTimeByteEq(data[0], 2))
	p.y.CondNegate(p.y.IsEven() ^ yShouldBeEven)
	p.z.SetUint64(1)
	p.normalized = true
	return nil
}

func (*Point) Curve() kyokusen.Curve {
	// TODO: Implement
	return nil
}

func (p1 *Point) Add(other kyokusen.Point) kyokusen.Point {
	p2 := castPoint(other)

	// This formula is taken from Algorithm 7 of https://eprint.iacr.org/2015/1060.
	t0 := NewField().Set(p1.x).Mul(p2.x)
	t1 := NewField().Set(p1.y).Mul(p2.y)
	t2 := NewField().Set(p1.z).Mul(p2.z)

	t3 := NewField().Set(p1.x).Add(p1.y)
	t4 := NewField().Set(p2.x).Add(p2.y)
	t3.Mul(t4)

	t4.Set(t0).Add(t1)
	t3.Sub(t4)
	t4.Set(p1.y).Add(p1.z)

	x := NewField().Set(p2.y).Add(p2.z)
	t4.Mul(x)
	x.Set(t1).Add(t2)

	t4.Sub(x)
	x.Set(p1.x).Add(p1.z)
	y := NewField().Set(p2.x).Add(p2.z)

	x.Mul(y)
	y.Set(t0).Add(t2)
	y.Negate().Add(x)

	x.Set(t0).Add(t0)
	t0.Add(x)
	t2.MulU64(3 * b)

	z := NewField().Set(t1).Add(t2)
	t1.Sub(t2)
	y.MulU64(3 * b)

	x.Set(y).Mul(t4)
	t2.Set(t3).Mul(t1)
	x.Negate().Add(t2)

	y.Mul(t0)
	t1.Mul(z)
	y.Add(t1)

	t0.Mul(t3)
	z.Mul(t4)
	z.Add(t0)

	return &Point{x, y, z, false}
}

func (p *Point) Sub(other kyokusen.Point) kyokusen.Point {
	return p.Add(other.Negate())
}

func (p *Point) Negate() kyokusen.Point {
	return &Point{
		x: NewField().Set(p.x),
		y: NewField().Set(p.y).Negate(),
		z: NewField().Set(p.z),
	}
}

func (p1 *Point) Equal(other kyokusen.Point) bool {
	p2 := castPoint(other)
	p1.normalize()
	p2.normalize()
	return (p1.x.Eq(p2.x) & p1.y.Eq(p2.y) & p1.z.Eq(p2.z)) == 1
}

func (p *Point) IsIdentity() bool {
	// Whenever Z == 0, this is the point at infinity.
	return p.z.EqZero() == 1
}

func (*Point) XScalar() kyokusen.Scalar {
	// TODO: Implement
	return nil
}

// CondAssign conditionally modifies the contents of a point.
func (p *Point) CondAssign(yes saferith.Choice, other *Point) *Point {
	p.x.CondAssign(yes, other.x)
	p.y.CondAssign(yes, other.y)
	p.z.CondAssign(yes, other.z)
	return p
}
