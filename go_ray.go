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

type point vector

type camera struct {
	eye point
	look_at point
	up vector
}

type screen struct {
	w,h float64
	xres, yres int
}

type scene struct {
	center point
	radius float64
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
	fmt.Println("hi")
	g_screen := screen{100,100,1000,1000}

	f, err := os.OpenFile("x.png", os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m := image.NewRGBA(image.Rect(0,0,g_screen.xres,g_screen.yres))
	
	for i:=0; i < g_screen.xres; i++ {
		for j:=0; j < g_screen.xres; j++ {
			m.Set(i,j,color.RGBA{255,0,0,255})
		}
	}
	if err=png.Encode(f,m); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}