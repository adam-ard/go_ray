package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"
)

const PI = 3.14159

var g_scene_desc scene_desc
var the_scene *scene
var translate_val float32
var rotate_val float32

type ray struct {
	start     vector
	direction vector
}

type scene struct {
	items []sceneItem `json:"items"`
}

type sceneItem interface {
	intersected(c_ray *ray) (float64, bool) // returns the t for the intersection, if it occured
	getReflectiveness() float64
	getColorRaw() (float64, float64, float64)
	getUnitNormal(point *vector) *vector
}

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {
	translate_val = 7.0
	rotate_val = 5.0
	scene, err := ReadFile("scene.json")
	if err != nil {
		fmt.Printf("Problem reading scene file: %s\n", err.Error())
		os.Exit(-1)
	}

	g_scene_desc = scene_desc{}
	if err := json.Unmarshal(scene, &g_scene_desc); err != nil {
		fmt.Printf("Problem parsing scene file: %s\n", err.Error())
		os.Exit(-1)
	}

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(800, 600, "Raytracer", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	window.SetKeyCallback(keyboard_callback)

	setupScene()
	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func keyboard_callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//	fmt.Println(key, scancode, action, mods)
	if key == glfw.KeyT && action == glfw.Press {
		translate_val = translate_val + 0.1
	}
	if key == glfw.KeyR && action == glfw.Repeat {
		rotate_val = rotate_val + 1.0
	}
	if key == glfw.KeyX && action == glfw.Repeat {
		g_scene_desc.Spheres[0].Center.X = g_scene_desc.Spheres[0].Center.X + 0.1
	}
	if key == glfw.KeyY && action == glfw.Repeat {
		g_scene_desc.Spheres[0].Center.Y = g_scene_desc.Spheres[0].Center.Y + 0.1
	}
	if key == glfw.KeyD && action == glfw.Press {
		trace("x.png")
	}
}

func setupScene() {
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.LIGHTING)
	gl.ColorMaterial(gl.FRONT_AND_BACK, gl.AMBIENT_AND_DIFFUSE)
	gl.Enable(gl.COLOR_MATERIAL)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)

	ambient_val := float32(g_scene_desc.AmbientLight)
	ambient := []float32{ambient_val, ambient_val, ambient_val, 1}
	diffuse := []float32{1, 1, 1, 1}

	light_pos := g_scene_desc.Light
	lightPosition := []float32{float32(light_pos.X), float32(light_pos.Y), float32(light_pos.Z), 0}
	gl.Lightfv(gl.LIGHT0, gl.AMBIENT, &ambient[0])
	gl.Lightfv(gl.LIGHT0, gl.DIFFUSE, &diffuse[0])
	gl.Lightfv(gl.LIGHT0, gl.POSITION, &lightPosition[0])
	gl.Enable(gl.LIGHT0)

	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1, 1, -1, 1, 1.0, 10.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func drawSphere(radius float32, lats, longs int) {
	for i := 0; i <= lats; i++ {
		lat0 := PI * (-0.5 + (float32(i) - 1.0/float32(lats)))
		z0 := math.Sin(float64(lat0))
		zr0 := math.Cos(float64(lat0))

		lat1 := PI * (-0.5 + (float32(i) / float32(lats)))
		z1 := math.Sin(float64(lat1))
		zr1 := math.Cos(float64(lat1))

		gl.Begin(gl.QUAD_STRIP)
		for j := 0; j <= longs; j++ {
			lng := 2.0 * PI * (float32(j) - 1.0) / float32(longs)
			x := math.Cos(float64(lng))
			y := math.Sin(float64(lng))

			gl.Normal3f(float32(x*zr0), float32(y*zr0), float32(z0))
			gl.Vertex3f(float32(x*zr0), float32(y*zr0), float32(z0))
			gl.Normal3f(float32(x*zr1), float32(y*zr1), float32(z1))
			gl.Vertex3f(float32(x*zr1), float32(y*zr1), float32(z1))
		}
		gl.End()
	}
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	gl.Translatef(0, 0, -translate_val)
	gl.Rotatef(rotate_val, 1, 0, 0)
	gl.Rotatef(rotate_val, 0, 1, 0)
	for _, v := range g_scene_desc.Spheres {

		gl.Translatef(float32(v.Center.X), float32(v.Center.Y), float32(v.Center.Z))
		gl.Color3f(float32(v.Red), float32(v.Green), float32(v.Blue))

		drawSphere(0.5, 200, 200)

		//		gl.Begin(gl.QUADS)
		//		gl.Normal3f(0, 0, 1)
		//		gl.Vertex3f(-1, -1, 1)
		//		gl.Vertex3f(1, -1, 1)
		//		gl.Vertex3f(1, 1, 1)
		//		gl.Vertex3f(-1, 1, 1)
		//
		//		gl.Normal3f(0, 0, -1)
		//		gl.Vertex3f(-1, -1, -1)
		//		gl.Vertex3f(-1, 1, -1)
		//		gl.Vertex3f(1, 1, -1)
		//		gl.Vertex3f(1, -1, -1)
		//
		//		gl.Normal3f(0, 1, 0)
		//		gl.Vertex3f(-1, 1, -1)
		//		gl.Vertex3f(-1, 1, 1)
		//		gl.Vertex3f(1, 1, 1)
		//		gl.Vertex3f(1, 1, -1)
		//
		//		gl.Normal3f(0, -1, 0)
		//		gl.Vertex3f(-1, -1, -1)
		//		gl.Vertex3f(1, -1, -1)
		//		gl.Vertex3f(1, -1, 1)
		//		gl.Vertex3f(-1, -1, 1)

		//		gl.Normal3f(1, 0, 0)
		//		gl.Vertex3f(1, -1, -1)
		//		gl.Vertex3f(1, 1, -1)
		//		gl.Vertex3f(1, 1, 1)
		//		gl.Vertex3f(1, -1, 1)
		//
		//		gl.Normal3f(-1, 0, 0)
		//		gl.Vertex3f(-1, -1, -1)
		//		gl.Vertex3f(-1, -1, 1)
		//		gl.Vertex3f(-1, 1, 1)
		//		gl.Vertex3f(-1, 1, -1)
		//		gl.End()
		gl.Translatef(float32(-v.Center.X), float32(-v.Center.Y), float32(-v.Center.Z))
	}

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
	the_scene := scene{}
	for _, v := range g_scene_desc.YPlanes {
		new_v := v
		the_scene.items = append(the_scene.items, &new_v)
	}
	for _, v := range g_scene_desc.Spheres {
		new_v := v
		the_scene.items = append(the_scene.items, &new_v)
	}
	return &the_scene
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

//func main() {
//	trace("scene.png", "x.png")
//}

func trace(targetFile string) {
	f, err := os.OpenFile(targetFile, os.O_CREATE|os.O_WRONLY, 0666)
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
