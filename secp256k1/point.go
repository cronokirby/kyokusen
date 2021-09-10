package secp256k1

// Point represents a point on the secp256k1 curve.
type Point struct {
	// Internally, we represent this as a projective point (X : Y : Z).
	// This corresponds to the affine point (X / Z, Y / Z), except when Z = 0,
	// which corresponds to the point at infinity.
	x *Field
	y *Field
	z *Field
}
