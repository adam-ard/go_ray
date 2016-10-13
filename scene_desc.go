package main

type camera struct {
	Eye     vector `json:"eye"`
	Look_at vector `json:"look_at"`
	Up      vector `json:"up"`
}

type screen struct {
	W    float64 `json:"w"`
	H    float64 `json:"h"`
	Xres int     `json:"xres"`
	Yres int     `json:"yres"`
}

type scene_desc struct {
	Screen       screen  `json:"screen"`
	Camera       camera  `json:"camera"`
	Light        vector  `json:"light"`
	AmbientLight float64 `json:"ambient_light"`
}
