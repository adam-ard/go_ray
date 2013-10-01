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
	radius float64
	red,blue,green uint8
}

func (s *sphere) getColor(c_ray *ray) (uint8, uint8, uint8, float64, bool) {
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

	// test with a sphere
	var red,blue,green uint8
	is_hit:=false
	i_test:=b*b-4.0*a*c
	t1,t2,t_closest:=0.0,0.0,0.0
	if i_test > 0.0 {
		red,blue,green=s.red,s.blue,s.green
		is_hit=true
		t1=(-b+math.Sqrt(i_test))/(2.0*a)
		t2=(-b-math.Sqrt(i_test))/(2.0*a)
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
	return red,blue,green,t_closest,is_hit
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

func main() {
	g_screen := screen{100,100,1000,1000}
	g_camera := camera{vector{0,0,100}, vector{0,0,0}, vector{0,1,0}}

	f, err := os.OpenFile("x.png", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m := image.NewRGBA(image.Rect(0,0,g_screen.xres,g_screen.yres))
	
	for i:=0; i < g_screen.xres; i++ {
		for j:=0; j < g_screen.yres; j++ {
			a_to_e := g_camera.eye.sub(&g_camera.look_at)
			w := a_to_e.unit()

			up_X_a_to_e := g_camera.up.cross(&a_to_e)
			u := up_X_a_to_e.unit()
			v := w.cross(&u)
			
			cu := (((2.0*float64(i) + 1.0)/(2.0 * float64(g_screen.xres))) - 0.5) * g_screen.w
			cv := (((2.0*float64(j) + 1.0)/(2.0 * float64(g_screen.yres))) - 0.5) * g_screen.h
				
			ucu := u.scalarMult(cu)
			ucv := v.scalarMult(cv)
			a_plus_ucu := g_camera.look_at.add(&ucu)
			Pij := a_plus_ucu.add(&ucv)

			e_to_Pij := Pij.sub(&g_camera.eye)

			current_ray := ray{g_camera.eye, e_to_Pij.unit()}
			
			// put the sphere at the origin
			s := sphere{vector{5.0, 15.0, 0.0}, 5.0, 0, 0, 255}
			s2 := sphere{vector{0.0, 0.0, -15.0}, 15.0, 255, 255, 0}
			
			red, blue, green, t, is_hit := s.getColor(&current_ray)
			red2, blue2, green2, t2, is_hit2 := s2.getColor(&current_ray)
			
			if is_hit && is_hit2 {
				if t<t2 {
					m.Set(i,j,color.RGBA{red,blue,green,255})
				}else{
					m.Set(i,j,color.RGBA{red2,blue2,green2,255})
				}
			} else if is_hit {
				m.Set(i,j,color.RGBA{red,blue,green,255})
			} else if is_hit2 {
				m.Set(i,j,color.RGBA{red2,blue2,green2,255})
			} else {
				m.Set(i,j,color.RGBA{0,0,0,255})
			}
		}
	}
	if err=png.Encode(f,m); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}