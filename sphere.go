package main

import "math"

type sphere struct {
	Center         vector  `json:"center"`
	Radius         float64 `json:"radius"`
	Reflectiveness float64 `json:"reflectiveness"`
	Red            float64 `json:"red"`
	Green          float64 `json:"green"`
	Blue           float64 `json:"blue"`
}

func (s *sphere) getReflectiveness() float64 {
	return s.Reflectiveness
}

func (s *sphere) getColorRaw() (float64, float64, float64) {
	return s.Red, s.Green, s.Blue
}

func (s *sphere) getUnitNormal(point *vector) *vector {
	normal := point.sub(&s.Center)
	unormal := normal.unit()
	return &unormal
}

func (s *sphere) intersected(c_ray *ray) (float64, bool) {
	a := c_ray.direction.X*c_ray.direction.X +
		c_ray.direction.Y*c_ray.direction.Y +
		c_ray.direction.Z*c_ray.direction.Z
	b := 2.0 * ((c_ray.start.X-s.Center.X)*c_ray.direction.X +
		(c_ray.start.Y-s.Center.Y)*c_ray.direction.Y +
		(c_ray.start.Z-s.Center.Z)*c_ray.direction.Z)
	c := (c_ray.start.X-s.Center.X)*(c_ray.start.X-s.Center.X) +
		(c_ray.start.Y-s.Center.Y)*(c_ray.start.Y-s.Center.Y) +
		(c_ray.start.Z-s.Center.Z)*(c_ray.start.Z-s.Center.Z) -
		s.Radius*s.Radius

	is_hit := false
	i_test := b*b - 4.0*a*c
	t1, t2, t_closest := 0.0, 0.0, 0.0
	if i_test >= 0.0 {
		is_hit = true
		t1 = (-b + math.Sqrt(i_test)) / (2.0 * a)
		t2 = (-b - math.Sqrt(i_test)) / (2.0 * a)
		t1 = in_buffer(t1)
		t2 = in_buffer(t2)
		if t1 <= 0.0 && t2 <= 0.0 {
			is_hit = false // it hit behind or on the viewer
		} else if t1 > 0.0 && t2 > 0.0 {
			if t1 < t2 {
				t_closest = t1
			} else {
				t_closest = t2
			}
		} else if t1 > 0.0 {
			t_closest = t1
		} else if t2 > 0.0 {
			t_closest = t2
		}
	}

	return t_closest, is_hit
}
