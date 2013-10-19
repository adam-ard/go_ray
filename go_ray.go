package main

import "math"
import "fmt"
import "image"
import "os"
import "image/color"
import "image/png"

// vector can represent a point as well, if interpreted as
//  a vector from the origin to a point or 'absolute position vector'.
//  otherwise it represents a 'relative position vector'
type vector struct {
	x, y, z float64
}

type ray struct {
	start vector
	direction vector
}

type camera struct {
	eye vector
	look_at vector
	up vector
}

type screen struct {
	w,h float64
	xres, yres int
}

type sphere struct {
	center vector
	radius,reflectiveness float64
	red,green,blue uint16
}

type scene struct {
	items []sceneItem
}

type zplane struct {
	loc, reflectiveness float64
	red,green,blue uint16
}

type sceneItem interface {
	intersected(c_ray *ray) (float64, bool)   // returns the t for the intersection, if it occured
	getColor(the_scene *scene, t, ambient float64, c_ray *ray, light *vector) (uint16, uint16, uint16) // get the color at intersection point
}

func (z *zplane) getReflectiveness() (float64) {
	return z.reflectiveness
}

func (s *sphere) getReflectiveness() (float64) {
	return s.reflectiveness
}

func (z *zplane) getColorRaw() (uint16, uint16, uint16) {
	return z.red, z.green, z.blue
}

func (s *sphere) getColorRaw() (uint16, uint16, uint16) {
	return s.red, s.green, s.blue
}

func (z *zplane) getUnitNormal(point *vector) (*vector) {
	return &vector{0.0, 0.0, -1.0}
}

func (s *sphere) getUnitNormal(point *vector) (*vector) {
	normal := point.sub(&s.center)
	unormal := normal.unit()
	return &unormal
}

func (the_scene *scene) getColor(c_ray *ray, light *vector, ambient float64) (uint16, uint16, uint16) {
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
		return 20000,20000,20000
	}
	
	return closest_item.getColor(the_scene, t_closest, ambient, c_ray, light)
}

func (z *zplane) intersected(c_ray *ray) (float64, bool)  {
	if c_ray.direction.z == 0.0 {
		return 0.0, false
	}

	t := (z.loc - c_ray.start.z) / c_ray.direction.z
	t = in_buffer(t)
	if t <= 0.0 {
		return 0.0, false
	}

	return t, true
}

func (z *zplane) getColor(the_scene *scene, t, ambient float64, c_ray *ray, light *vector) (uint16, uint16, uint16) {
	dir := c_ray.direction.scalarMult(t)
	point_on_plane:= c_ray.start.add(&dir)
	
	unormal := z.getUnitNormal(&point_on_plane)

	// check for light obstructions
	point_on_plane_to_light := light.sub(&point_on_plane)
	upoint_on_plane_to_light := point_on_plane_to_light.unit()
	
	is_hit:=false
	scale:=0.0
	for _, value := range the_scene.items {
		_, is_hit = value.intersected(&ray{point_on_plane, upoint_on_plane_to_light})
		if is_hit {
			break
		}
	}
	
	//calculate the light contribution
	upoint_on_plane_to_source := c_ray.direction.scalarMult(-1.0)
	intermediate := unormal.scalarMult(2.0 * upoint_on_plane_to_source.dot(unormal))
	reflected := intermediate.sub(&upoint_on_plane_to_source)
	ureflected := reflected.unit()
		
	if is_hit == false {
		scale = ureflected.dot(&upoint_on_plane_to_light)
		if scale < 0.0 {
			scale = 0.0
		}
	}

	red_light := scale*float64(z.red)
	green_light := scale*float64(z.green)
	blue_light := scale*float64(z.blue)
	
	obj_red:=ceiling(red_light + ambient*float64(z.red),65535)
	obj_green:=ceiling(green_light + ambient*float64(z.green),65535)
	obj_blue:=ceiling(blue_light + ambient*float64(z.blue),65535)

	// send the reflectived ray into the scene
	reflected_red,reflected_green,reflected_blue:=the_scene.getColor(&ray{point_on_plane, ureflected},light,ambient)

	return uint16(z.reflectiveness*float64(reflected_red) + (1.0-z.reflectiveness)*obj_red),
	uint16(z.reflectiveness*float64(reflected_green) + (1.0-z.reflectiveness)*obj_green), 
	uint16(z.reflectiveness*float64(reflected_blue) + (1.0-z.reflectiveness)*obj_blue)
}

func (s *sphere) getColor(the_scene *scene, t, ambient float64, c_ray *ray, light *vector) (uint16, uint16, uint16) {
	// get the normal
	dir := c_ray.direction.scalarMult(t)
	point_on_sphere := c_ray.start.add(&dir)
	
	unormal := s.getUnitNormal(&point_on_sphere)
	
	// check for light obstructions
	point_on_sphere_to_light := light.sub(&point_on_sphere)
	upoint_on_sphere_to_light := point_on_sphere_to_light.unit()
	
	is_hit:=false
	scale:=0.0
	for _, value := range the_scene.items {
		_, is_hit := value.intersected(&ray{point_on_sphere, upoint_on_sphere_to_light})
		if is_hit {
			break
		}
	}
	
	//calculate the light contribution
	upoint_on_sphere_to_source := c_ray.direction.scalarMult(-1.0)
	intermediate := unormal.scalarMult(2.0 * upoint_on_sphere_to_source.dot(unormal))
	reflected := intermediate.sub(&upoint_on_sphere_to_source)
	ureflected := reflected.unit()
	
	if is_hit == false {
		scale := ureflected.dot(&upoint_on_sphere_to_light)
		if scale < 0.0 {
			scale = 0.0
		}
	}

	red_light := scale*float64(s.red)
	green_light := scale*float64(s.green)
	blue_light := scale*float64(s.blue)

	obj_red:=ceiling(red_light + ambient*float64(s.red),65535)
	obj_green:=ceiling(green_light + ambient*float64(s.green),65535)
	obj_blue:=ceiling(blue_light + ambient*float64(s.blue),65535)

	// send the reflectived ray into the scene
	reflected_red,reflected_green,reflected_blue:=the_scene.getColor(&ray{point_on_sphere,ureflected},light,ambient)

	return uint16(s.reflectiveness*float64(reflected_red) + (1.0-s.reflectiveness)*obj_red), 
	uint16(s.reflectiveness*float64(reflected_green) + (1.0-s.reflectiveness)*obj_green), 
	uint16(s.reflectiveness*float64(reflected_blue) + (1.0-s.reflectiveness)*obj_blue)
}

var buffer_val float64 = .00001
func in_buffer(val float64) float64{
	if val < buffer_val {
		val = 0.0
	}
	return val
}

func (s *sphere) intersected(c_ray *ray) (float64, bool)  {
	a := c_ray.direction.x * c_ray.direction.x +
		c_ray.direction.y * c_ray.direction.y +
		c_ray.direction.z * c_ray.direction.z
	b := 2.0*((c_ray.start.x-s.center.x)*c_ray.direction.x +
		(c_ray.start.y-s.center.y)*c_ray.direction.y +
		(c_ray.start.z-s.center.z)*c_ray.direction.z) 
	c := (c_ray.start.x-s.center.x) * (c_ray.start.x-s.center.x) +
		(c_ray.start.y-s.center.y) * (c_ray.start.y-s.center.y) +
		(c_ray.start.z-s.center.z) * (c_ray.start.z-s.center.z) -
		s.radius * s.radius

	is_hit:=false
	i_test:=b*b-4.0*a*c
	t1,t2,t_closest:=0.0,0.0,0.0
	if i_test >= 0.0 {
		is_hit=true
		t1=(-b+math.Sqrt(i_test))/(2.0*a)
		t2=(-b-math.Sqrt(i_test))/(2.0*a)
		t1=in_buffer(t1)
		t2=in_buffer(t2)
		if t1 <= 0.0 && t2 <= 0.0 {
			is_hit=false  // it hit behind or on the viewer
		}else if t1 > 0.0 && t2 > 0.0 {
			if t1<t2 {
				t_closest=t1
			}else{
				t_closest=t2
			}
		} else if t1 > 0.0 {
			t_closest = t1
		} else if t2 > 0.0 {
			t_closest = t2
		}
	}
	
	return t_closest, is_hit
}

func (v *vector) sub(v1 *vector) vector {
	return vector{v.x-v1.x, v.y-v1.y, v.z-v1.z}
} 

func (v *vector) add(v1 *vector) vector {
	return vector{v.x+v1.x, v.y+v1.y, v.z+v1.z}
}

func (v *vector) scalarMult (c float64) vector {
	return vector{c * v.x, c * v.y, c * v.z}
}

func (v *vector) lengthSq() float64 {
	return v.x * v.x  + v.y * v.y + v.z * v.z
}

func (v *vector) length() float64 {
	return math.Sqrt(v.lengthSq())
}

func (v *vector) unit() vector {
	l := v.length()
	return vector{v.x/l, v.y/l, v.z/l}
}

func (v *vector) dot(v1 *vector) float64 {
	return v.x * v1.x + v.y * v1.y + v.z * v1.z
}

func (v1 *vector) cross(v2 *vector) vector {
	return vector{v1.y * v2.z - v2.y * v1.z,
		v2.x * v1.z - v1.x * v2.z,
		v1.x * v2.y - v2.x * v1.y}
}

func ceiling(value float64, top_value float64) float64{
	if value > top_value {
		value=top_value
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

func get_scene() (*scene) {
	s := sphere{vector{-25.0, 15.0, -20.0}, 10.0, 0.25, 0, 0, 65535}
	s2 := sphere{vector{-5.0, -15.0, -15.0}, 15.0, 0.25, 0, 65535, 0}
	s3 := sphere{vector{5.0, 15.0, -15.0}, 15.0, 0.25, 65535, 0, 0}
	z := zplane{-30.0, 0.70, 65535, 65535, 65535}
	the_scene:=new(scene)
	the_scene.items=make([]sceneItem,4)
	the_scene.items[0]=&s
	the_scene.items[1]=&s2
	the_scene.items[2]=&s3
	the_scene.items[3]=&z
	return the_scene
}

func get_current_ray (i, j int, the_screen *screen, u, v, look_at, eye *vector) (*ray) {
	cu := (((2.0*float64(i) + 1.0)/(2.0 * float64(the_screen.xres))) - 0.5) * the_screen.w
	cv := (((2.0*float64(j) + 1.0)/(2.0 * float64(the_screen.yres))) - 0.5) * the_screen.h
	ucu := u.scalarMult(cu)
	ucv := v.scalarMult(cv)
	a_plus_ucu := look_at.add(&ucu)
	Pij := a_plus_ucu.add(&ucv)
	e_to_Pij := Pij.sub(eye)
	current_ray := ray{*eye, e_to_Pij.unit()}
	return &current_ray
}

func main() {
	g_screen := screen{100,100,1000,1000}
	g_camera := camera{vector{-100,-100,5}, vector{0,0,0}, vector{0,0,-1}}
	g_light := vector{0.0,100.0,100.0}
	g_ambient := 0.2

	f, err := os.OpenFile("x.png", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m := image.NewRGBA64(image.Rect(0,0,g_screen.xres,g_screen.yres))
	the_scene:=get_scene()
	u, v := get_local_coordinate_system(&g_camera.eye, &g_camera.look_at, &g_camera.up)
	for i:=0; i < g_screen.xres; i++ {
		for j:=0; j < g_screen.yres; j++ {
			current_ray:=get_current_ray(i, j, &g_screen, u, v, &g_camera.look_at, &g_camera.eye)
			// shoot the ray into the scene
			red,green,blue:=the_scene.getColor(current_ray, &g_light, g_ambient)
			m.Set(i,j,color.RGBA64{red,green,blue,65535})
		}
	}
	if err=png.Encode(f,m); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}