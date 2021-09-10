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

func (*Point) Equal(kyokusen.Point) bool {
	// TODO: Implement
	return false
}

func (*Point) IsIdentity() bool {
	// TODO: Implement
	return false
}

func (*Point) XScalar() kyokusen.Scalar {
	// TODO: Implement
	return nil
}
