package main

import "math"

// vector can represent a point as well, if interpreted as
//  a vector from the origin to a point or 'absolute position vector'.
//  otherwise it represents a 'relative position vector'
type vector struct {
	x, y, z float64
}

func (v *vector) sub(v1 *vector) vector {
	return vector{v.x - v1.x, v.y - v1.y, v.z - v1.z}
}

func (v *vector) add(v1 *vector) vector {
	return vector{v.x + v1.x, v.y + v1.y, v.z + v1.z}
}

func (v *vector) scalarMult(c float64) vector {
	return vector{c * v.x, c * v.y, c * v.z}
}

func (v *vector) lengthSq() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

func (v *vector) length() float64 {
	return math.Sqrt(v.lengthSq())
}

func (v *vector) unit() vector {
	l := v.length()
	return vector{v.x / l, v.y / l, v.z / l}
}

func (v *vector) dot(v1 *vector) float64 {
	return v.x*v1.x + v.y*v1.y + v.z*v1.z
}

func (v1 *vector) cross(v2 *vector) vector {
	return vector{v1.y*v2.z - v2.y*v1.z,
		v2.x*v1.z - v1.x*v2.z,
		v1.x*v2.y - v2.x*v1.y}
}
