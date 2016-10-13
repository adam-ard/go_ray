package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
)

type ray struct {
	start     vector
	direction vector
}

type scene struct {
	items []sceneItem
}

type sceneItem interface {
	intersected(c_ray *ray) (float64, bool) // returns the t for the intersection, if it occured
	getReflectiveness() float64
	getColorRaw() (float64, float64, float64)
	getUnitNormal(point *vector) *vector
}

func (the_scene *scene) getColor(c_ray *ray, light *vector, ambient float64) (float64, float64, float64) {
	t_closest := 0.0
	var closest_item sceneItem = nil
	for _, value := range the_scene.items {
		t, is_hit := value.intersected(c_ray)
		if is_hit == false {
			continue
		}
		if closest_item == nil || (t > 0 && t < t_closest) {
			t_closest = t
			closest_item = value
		}
	}

	if closest_item == nil {
		return 0, 0, 0
	}

	dir := c_ray.direction.scalarMult(t_closest)
	point_on_object := c_ray.start.add(&dir)

	unormal := closest_item.getUnitNormal(&point_on_object)

	// check for light obstructions
	point_on_object_to_light := light.sub(&point_on_object)
	upoint_on_object_to_light := point_on_object_to_light.unit()

	is_hit := false
	scale := 0.0
	for _, value := range the_scene.items {
		_, is_hit = value.intersected(&ray{point_on_object, upoint_on_object_to_light})
		if is_hit {
			break
		}
	}

	//calculate the light contribution
	upoint_on_object_to_source := c_ray.direction.scalarMult(-1.0)
	intermediate := unormal.scalarMult(2.0 * upoint_on_object_to_source.dot(unormal))
	reflected := intermediate.sub(&upoint_on_object_to_source)
	ureflected := reflected.unit()

	if is_hit == false {
		scale = ureflected.dot(&upoint_on_object_to_light)
		if scale < 0.0 {
			scale = 0.0
		}
	}

	red, green, blue := closest_item.getColorRaw()

	red_light := scale * red
	green_light := scale * green
	blue_light := scale * blue

	obj_red := ceiling(red_light+ambient*red, 1.0)
	obj_green := ceiling(green_light+ambient*green, 1.0)
	obj_blue := ceiling(blue_light+ambient*blue, 1.0)

	// send the reflectived ray into the scene
	reflected_red, reflected_green, reflected_blue := the_scene.getColor(&ray{point_on_object, ureflected}, light, ambient)

	reflectiveness := closest_item.getReflectiveness()

	return reflectiveness*reflected_red + (1.0-reflectiveness)*obj_red,
		reflectiveness*reflected_green + (1.0-reflectiveness)*obj_green,
		reflectiveness*reflected_blue + (1.0-reflectiveness)*obj_blue
}

var buffer_val float64 = .00001

func in_buffer(val float64) float64 {
	if val < buffer_val {
		val = 0.0
	}
	return val
}

func ceiling(value float64, top_value float64) float64 {
	if value > top_value {
		value = top_value
	}
	return value
}

func get_local_coordinate_system(eye, look_at, up *vector) (*vector, *vector) {
	a_to_e := eye.sub(look_at)
	w := a_to_e.unit()
	up_X_a_to_e := up.cross(&a_to_e)
	u := up_X_a_to_e.unit()
	v := w.cross(&u)
	return &u, &v
}

func get_scene() *scene {
	s := sphere{vector{-25.0, 10.0, -20.0}, 10.0, 0.25, 0, 0, 0.75}
	s2 := sphere{vector{5.0, 15.0, 15.0}, 15.0, 0.25, 0, 0.75, 0}
	s3 := sphere{vector{5.0, 15.0, -15.0}, 15.0, 0.25, 0.75, 0, 0}
	y := yplane{0.0, 0.70, 1.0, 1.0, 1.0}
	the_scene := new(scene)
	the_scene.items = make([]sceneItem, 4)
	the_scene.items[0] = &s
	the_scene.items[1] = &s2
	the_scene.items[2] = &s3
	the_scene.items[3] = &y
	return the_scene
}

func get_current_ray(i, j int, the_screen *screen, u, v, look_at, eye *vector) *ray {
	cu := (((2.0*float64(i) + 1.0) / (2.0 * float64(the_screen.Xres))) - 0.5) * the_screen.W
	cv := (((2.0*float64(j) + 1.0) / (2.0 * float64(the_screen.Yres))) - 0.5) * the_screen.H
	ucu := u.scalarMult(cu)
	ucv := v.scalarMult(cv)
	a_plus_ucu := look_at.add(&ucu)
	Pij := a_plus_ucu.add(&ucv)
	e_to_Pij := Pij.sub(eye)
	current_ray := ray{*eye, e_to_Pij.unit()}
	return &current_ray
}

func ReadFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func main() {
	scene, err := ReadFile("scene.json")
	if err != nil {
		fmt.Printf("Problem reading scene file: %s", err.Error())
		os.Exit(-1)
	}

	g_scene_desc := scene_desc{}
	if err := json.Unmarshal(scene, &g_scene_desc); err != nil {
		fmt.Printf("Problem parsing scene file: %s", err.Error())
		os.Exit(-1)
	}

	fmt.Println(g_scene_desc)

	//g_scene_desc.Screen = screen{100, 100, 1000, 1000}
	//g_scene_desc.Camera = camera{vector{-100, 50, 50}, vector{5, 15, -15}, vector{0, 1, 0}}
	//g_scene_desc.Light = vector{0.0, 100.0, 100.0}
	//g_scene_desc.AmbientLight = 0.2

	f, err := os.OpenFile("x.png", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m := image.NewRGBA64(image.Rect(0, 0, g_scene_desc.Screen.Xres, g_scene_desc.Screen.Yres))
	the_scene := get_scene()
	u, v := get_local_coordinate_system(&g_scene_desc.Camera.Eye, &g_scene_desc.Camera.Look_at, &g_scene_desc.Camera.Up)

	num_segments := 1
	d := make(chan bool, num_segments)
	step := g_scene_desc.Screen.Xres / num_segments
	for segment := 0; segment < num_segments; segment++ {
		go func(c_seg int) {
			//			fmt.Println(c_seg, step)
			for i := c_seg * step; i < (c_seg+1)*step; i++ {
				for j := 0; j < g_scene_desc.Screen.Yres; j++ {
					current_ray := get_current_ray(i, j, &g_scene_desc.Screen, u, v, &g_scene_desc.Camera.Look_at, &g_scene_desc.Camera.Eye)
					red, green, blue := the_scene.getColor(current_ray, &g_scene_desc.Light, g_scene_desc.AmbientLight)
					m.Set(i, g_scene_desc.Screen.Yres-j, color.RGBA64{uint16(65535 * red), uint16(65535 * green), uint16(65535 * blue), 65535})
				}
			}
			d <- true
		}(segment)
	}

	for count := 0; count < num_segments; count++ {
		<-d
	}

	if err = png.Encode(f, m); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
