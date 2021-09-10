package secp256k1

import "github.com/cronokirby/kyokusen"

// Point represents a point on the secp256k1 curve.
type Point struct {
	// Internally, we represent this as a projective point (X : Y : Z).
	// This corresponds to the affine point (X / Z, Y / Z), except when Z = 0,
	// which corresponds to the point at infinity.
	x *Field
	y *Field
	z *Field
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
}

func (*Point) MarshalBinary() ([]byte, error) {
	// TODO: Implement
	return nil, nil
}

func (*Point) UnmarshalBinary(data []byte) error {
	// TODO: Implement
	return nil
}

func (*Point) Curve() kyokusen.Curve {
	// TODO: Implement
	return nil
}

func (*Point) Add(kyokusen.Point) kyokusen.Point {
	// TODO: Implement
	return nil
}

func (*Point) Sub(kyokusen.Point) kyokusen.Point {
	// TODO: Implement
	return nil
}

func (*Point) Negate() kyokusen.Point {
	// TODO: Implement
	return nil
}

func (p1 *Point) Equal(other kyokusen.Point) bool {
	p2 := castPoint(other)
	p1.normalize()
	p2.normalize()
	return (p1.x.Eq(p2.x) & p1.y.Eq(p2.y) & p1.z.Eq(p2.z)) == 1
}

func (*Point) IsIdentity() bool {
	// TODO: Implement
	return false
}

func (*Point) XScalar() kyokusen.Scalar {
	// TODO: Implement
	return nil
}
