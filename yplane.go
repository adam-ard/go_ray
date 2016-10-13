package main

type yplane struct {
	Loc            float64 `json:"loc"`
	Reflectiveness float64 `json:"reflectiveness"`
	Red            float64 `json:"red"`
	Green          float64 `json:"green"`
	Blue           float64 `json:"blue"`
}

func (y *yplane) getReflectiveness() float64 {
	return y.Reflectiveness
}

func (y *yplane) getColorRaw() (float64, float64, float64) {
	return y.Red, y.Green, y.Blue
}

func (y *yplane) getUnitNormal(point *vector) *vector {
	return &vector{0.0, 1.0, 0.0}
}

func (y *yplane) intersected(c_ray *ray) (float64, bool) {
	if c_ray.direction.Y == 0.0 {
		return 0.0, false
	}

	t := (y.Loc - c_ray.start.Y) / c_ray.direction.Y
	t = in_buffer(t)
	if t <= 0.0 {
		return 0.0, false
	}

	return t, true
}
