package main

type yplane struct {
	loc, reflectiveness float64
	red, green, blue    float64
}

func (y *yplane) getReflectiveness() float64 {
	return y.reflectiveness
}

func (y *yplane) getColorRaw() (float64, float64, float64) {
	return y.red, y.green, y.blue
}

func (y *yplane) getUnitNormal(point *vector) *vector {
	return &vector{0.0, 1.0, 0.0}
}

func (y *yplane) intersected(c_ray *ray) (float64, bool) {
	if c_ray.direction.Y == 0.0 {
		return 0.0, false
	}

	t := (y.loc - c_ray.start.Y) / c_ray.direction.Y
	t = in_buffer(t)
	if t <= 0.0 {
		return 0.0, false
	}

	return t, true
}
