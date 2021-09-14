package secp256k1

import (
	"github.com/cronokirby/kyokusen"
	"github.com/cronokirby/saferith"
)

// Curve represents the secp256k1 curve, implementing the kyokusen.Curve interface.
type Curve struct{}

func (Curve) NewPoint() kyokusen.Point {
	return NewPoint()
}

func (Curve) NewBasePoint() kyokusen.Point {
	// Prevents potential concurrency issues with normalization. Internally, we don't
	// ever check for equality with the base point, so there's no issues with normalization
	// modifying the shared base point. We have no control over what users do.
	return NewPoint().CondAssign(1, basePoint)
}

func (Curve) NewScalar() kyokusen.Scalar {
	return NewScalar()
}

func (Curve) Name() string {
	return "secp256k1"
}

func (Curve) ScalarBits() int {
	return 256
}

func (Curve) SafeScalarBytes() int {
	return 32
}

func (Curve) Order() *saferith.Modulus {
	return q
}
