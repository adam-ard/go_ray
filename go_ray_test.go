package main

import "testing"
import "math"

func TestVectorSub(t *testing.T) {
	testTable := []struct {
		v1, v2, result vector
	}{
		{vector{2,3,4},vector{1,2,3},vector{1,1,1}},
		{vector{1,1,1},vector{1,1,1},vector{0,0,0}},
		{vector{0,0,0},vector{0,0,0},vector{0,0,0}},
		{vector{-1,-1,-1},vector{-1,-1,-1},vector{0,0,0}},
		{vector{1,1,1},vector{-100,-100,-100},vector{101,101,101}},
	}

	for _, value := range testTable {
		result := value.v1.sub(&value.v2)
		if  result != value.result {
			t.Error("ERROR: looking for", value.result, " got ", result)
		} 
	}
}

func TestVectorAdd(t *testing.T) {
	testTable := []struct {
		v1, v2, result vector
	}{
		{vector{2,3,4},vector{1,2,3},vector{3,5,7}},
		{vector{1,1,1},vector{1,1,1},vector{2,2,2}},
		{vector{0,0,0},vector{0,0,0},vector{0,0,0}},
		{vector{-1,-1,-1},vector{-1,-1,-1},vector{-2,-2,-2}},
		{vector{1,1,1},vector{-100,-100,-100},vector{-99,-99,-99}},
	}

	for _, value := range testTable {
		result := value.v1.add(&value.v2)
		if  result != value.result {
			t.Error("ERROR: looking for", value.result, " got ", result)
		} 
	}
}

func TestVectorDot(t *testing.T) {
	testTable := []struct {
		v1, v2 vector
		result float64
	}{
		{vector{2,3,4},vector{1,2,3},20},
		{vector{1,1,1},vector{1,1,1},3},
		{vector{0,0,0},vector{0,0,0},0},
		{vector{-1,-1,-1},vector{-1,-1,-1},3},
		{vector{1,1,1},vector{-100,-100,-100},-300},
	}

	for _, value := range testTable {
		result := value.v1.dot(&value.v2)
		if  result != value.result {
			t.Error("ERROR: looking for", value.result, " got ", result)
		} 
	}
}

func TestVectorScalarMult(t *testing.T) {
	testTable := []struct {
		v1 vector
		c float64
		result vector
	}{
		{vector{2,3,4},2,vector{4,6,8}},
		{vector{1,1,1},2,vector{2,2,2}},
		{vector{0,0,0},2,vector{0,0,0}},
		{vector{2,3,4},0,vector{0,0,0}},
		{vector{2,3,4},-2,vector{-4,-6,-8}},
		{vector{2,3,4},10000,vector{20000,30000,40000}},
	}

	for _, value := range testTable {
		result := value.v1.scalarMult(value.c)
		if  result != value.result {
			t.Error("ERROR: looking for", value.result, " got ", result)
		} 
	}
}

func TestVectorLengthSq(t *testing.T) {
	testTable := []struct {
		v vector
		result float64
	}{
		{vector{2,3,4},29},
		{vector{1,1,1},3},
		{vector{0,0,0},0},
	}

	for _, value := range testTable {
		result := value.v.lengthSq()
		if  result != value.result {
			t.Error("ERROR: looking for", value.result, " got ", result)
		} 
	}
}

func TestVectorLength(t *testing.T) {
	testTable := []struct {
		v vector
		result float64
	}{
		{vector{2,3,4},math.Sqrt(29)},
		{vector{1,1,1},math.Sqrt(3)},
		{vector{0,0,0},0},
	}

	for _, value := range testTable {
		result := value.v.length()
		if  result != value.result {
			t.Error("ERROR: looking for", value.result, " got ", result)
		} 
	}
}

func TestVectorCross(t *testing.T) {
	testTable := []struct {
		v1, v2, result vector
	}{
		{vector{2,3,4},vector{1,2,3},vector{1,-2,1}},
		{vector{1,1,1},vector{1,1,1},vector{0,0,0}},
		{vector{2,2,2},vector{1,1,1},vector{0,0,0}},
		{vector{-2,-2,-2},vector{1,1,1},vector{0,0,0}},
		{vector{-2,-3,-2},vector{1,1,1},vector{-1,0,1}},
	}

	for _, value := range testTable {
		result := value.v1.cross(&value.v2)
		if  result != value.result {
			t.Error("ERROR: looking for", value.result, " got ", result)
		} 
	}
}
