package main

import "math"

type sphere struct {
	center                 vector
	radius, reflectiveness float64
	red, green, blue       float64
}

func (s *sphere) getReflectiveness() float64 {
	return s.reflectiveness
}

func (s *sphere) getColorRaw() (float64, float64, float64) {
	return s.red, s.green, s.blue
}

func (s *sphere) getUnitNormal(point *vector) *vector {
	normal := point.sub(&s.center)
	unormal := normal.unit()
	return &unormal
}

func (s *sphere) intersected(c_ray *ray) (float64, bool) {
	a := c_ray.direction.X*c_ray.direction.X +
		c_ray.direction.Y*c_ray.direction.Y +
		c_ray.direction.Z*c_ray.direction.Z
	b := 2.0 * ((c_ray.start.X-s.center.X)*c_ray.direction.X +
		(c_ray.start.Y-s.center.Y)*c_ray.direction.Y +
		(c_ray.start.Z-s.center.Z)*c_ray.direction.Z)
	c := (c_ray.start.X-s.center.X)*(c_ray.start.X-s.center.X) +
		(c_ray.start.Y-s.center.Y)*(c_ray.start.Y-s.center.Y) +
		(c_ray.start.Z-s.center.Z)*(c_ray.start.Z-s.center.Z) -
		s.radius*s.radius

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
