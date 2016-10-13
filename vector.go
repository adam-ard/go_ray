package main

import "math"

// vector can represent a point as well, if interpreted as
//  a vector from the origin to a point or 'absolute position vector'.
//  otherwise it represents a 'relative position vector'
type vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (v *vector) sub(v1 *vector) vector {
	return vector{v.X - v1.X, v.Y - v1.Y, v.Z - v1.Z}
}

func (v *vector) add(v1 *vector) vector {
	return vector{v.X + v1.X, v.Y + v1.Y, v.Z + v1.Z}
}

func (v *vector) scalarMult(c float64) vector {
	return vector{c * v.X, c * v.Y, c * v.Z}
}

func (v *vector) lengthSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *vector) length() float64 {
	return math.Sqrt(v.lengthSq())
}

func (v *vector) unit() vector {
	l := v.length()
	return vector{v.X / l, v.Y / l, v.Z / l}
}

func (v *vector) dot(v1 *vector) float64 {
	return v.X*v1.X + v.Y*v1.Y + v.Z*v1.Z
}

func (v1 *vector) cross(v2 *vector) vector {
	return vector{v1.Y*v2.Z - v2.Y*v1.Z,
		v2.X*v1.Z - v1.X*v2.Z,
		v1.X*v2.Y - v2.X*v1.Y}
}
