package lmath

import (
	"fmt"
	"math"
	"testing"
)

// Test creation with 16 values
// Test Equals with matrices
func TestNewMat3(t *testing.T) {
	m := NewMat3(
		1, 2, 3,
		5, 6, 7,
		9, 10, 11)
	m2 := NewMat3(
		1, 2, 3,
		5, 6, 7,
		9, 10, 11)
	if m.Eq(*m2) == false {
		t.Error(fmt.Sprintf("TestNewMat3 "))
	}
}

// Test Get
// Test Set
func TestGetterSetterMat3(t *testing.T) {
	cases := []struct {
		orig          *Mat3
		x, y          int
		wantBeforeSet float64
		wantAfterSet  float64
	}{
		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 0, 0, 1, 10},
		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 1, 0, 2, 20},
		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 2, 0, 3, 30},

		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 0, 1, 4, 40},
		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 1, 1, 5, 50},
		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 2, 1, 6, 60},

		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 0, 2, 7, 70},
		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 1, 2, 8, 80},
		{&Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}, 2, 2, 9, 90},
	}

	for testIndex, c := range cases {
		get := c.orig.Get(c.y, c.x)
		if get != c.wantBeforeSet {
			t.Error(fmt.Sprintf("TestGetterSetterMat3 wantBeforeSet %d %v", testIndex, get))
		}

		get = c.orig.Set(c.y, c.x, c.wantAfterSet).Get(c.y, c.x)
		if get != c.wantAfterSet {
			t.Error(fmt.Sprintf("TestGetterSetterMat3 wantAfterSet %d %v", testIndex, get))
		}
	}

	orig := &Mat3{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}}
	for k, _ := range orig.Dump() {
		get := orig.At(k)
		if get != float64(k+1) {
			t.Error(fmt.Sprintf("TestGetterSetterMat3 At %d %v", k, get))
		}

		orig.SetAt(k, float64((k+1)*10))
		get = orig.At(k)
		if get != float64((k+1)*10) {
			t.Error(fmt.Sprintf("TestGetterSetterMat3 SetAt %d", k))
		}
	}
}

// Test Load Array
// Test Dump
// Test Dump
func TestLoadDumpMat3(t *testing.T) {
	cases := []struct {
		loadArray [9]float64
	}{
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{[9]float64{-1, -2, -3, -4, -5, -6, -7, -8, -9}},
	}

	m := &Mat3{}
	for testIndex, c := range cases {
		m.Load(c.loadArray)
		get := m.Dump()

		for k, _ := range get {
			if get[k] != c.loadArray[k] {
				t.Error(fmt.Sprintf("TestLoadDumpMat3 %d", testIndex))
				break
			}
		}
	}

	// 1 2 3
	// 4 5 6
	// 7 8 9
	m.Load([9]float64{1, 4, 7, 2, 5, 8, 3, 6, 9})
	get := m.DumpOpenGL()
	for k, _ := range get {
		if !closeEq(get[k], cases[0].loadArray[k], epsilon) {
			t.Error(fmt.Sprintf("TestDumpOpenGLMat3"))
			break
		}
	}

	get2 := m.DumpOpenGLf32()
	for k, _ := range get {
		if !closeEq(float64(get2[k]), cases[0].loadArray[k], epsilon) {
			t.Error(fmt.Sprintf("TestDumpOpenGLf32Mat3 %d", k))
			break
		}
	}
}

func TestRowMat3(t *testing.T) {
	cases := []struct {
		orig     [9]float64
		rowIndex int
		want     [3]float64
	}{
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0, [3]float64{1, 2, 3}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1, [3]float64{4, 5, 6}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2, [3]float64{7, 8, 9}},
	}

	m := &Mat3{}
	var x, y, z float64
	for testIndex, c := range cases {
		m.Load(c.orig)
		x, y, z = m.Row(c.rowIndex)
		if x != c.want[0] || y != c.want[1] || z != c.want[2] {
			t.Error(fmt.Sprintf("TestRowMat3 %d", testIndex))
		}
	}
}

func TestSetRowMat3(t *testing.T) {
	cases := []struct {
		orig     [9]float64
		rowIndex int
		x, y, z  float64
	}{
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0, -1, -2, -3},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1, -4, -5, -6},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2, -7, -8, -9},
	}

	m := &Mat3{}
	var x, y, z float64
	for testIndex, c := range cases {
		m.Load(c.orig)
		m.SetRow(c.rowIndex, c.x, c.y, c.z)
		x, y, z = m.Row(c.rowIndex)
		if x != c.x || y != c.y || z != c.z {
			t.Error(fmt.Sprintf("TestSetRowMat3 %d", testIndex))
		}
	}
}

func TestColMat3(t *testing.T) {
	cases := []struct {
		orig     [9]float64
		colIndex int
		want     [3]float64
	}{
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0, [3]float64{1, 4, 7}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1, [3]float64{2, 5, 8}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2, [3]float64{3, 6, 9}},
	}

	m := &Mat3{}
	var x, y, z float64
	for testIndex, c := range cases {
		m.Load(c.orig)
		x, y, z = m.Col(c.colIndex)
		if x != c.want[0] || y != c.want[1] || z != c.want[2] {
			t.Error(fmt.Sprintf("TestColMat3 %d", testIndex))
		}
	}
}

func TestSetColMat3(t *testing.T) {
	cases := []struct {
		orig     [9]float64
		colIndex int
		x, y, z  float64
	}{
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0, -1, -2, -3},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1, -5, -6, -7},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2, -9, -10, -11},
	}

	m := &Mat3{}
	var x, y, z float64
	for testIndex, c := range cases {
		m.Load(c.orig)
		m.SetCol(c.colIndex, c.x, c.y, c.z)
		x, y, z = m.Col(c.colIndex)
		if x != c.x || y != c.y || z != c.z {
			t.Error(fmt.Sprintf("TestSetColMat3 %d", testIndex))
		}
	}
}

func TestAddScalarMat3(t *testing.T) {
	cases := []struct {
		orig  [9]float64
		value float64
	}{
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, 1},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, -1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, -1},
	}

	m := Mat3{}
	for testIndex, c := range cases {
		m.Load(c.orig)
		ret_mat := m.AddScalar(c.value)

		get := ret_mat.Dump()
		for k, _ := range c.orig {
			if get[k] != c.orig[k]+c.value {
				t.Error(fmt.Sprintf("TestAddScalarMat3 %d %d", testIndex, k))
				break
			}
		}

		ret_mat2 := m.AddInScalar(c.value)
		if ret_mat2 != &m {
			t.Error(fmt.Sprintf("TestAddInScalarMat3 %d", testIndex))
		}

		get = ret_mat2.Dump()
		for k, _ := range c.orig {
			if get[k] != c.orig[k]+c.value {
				t.Error(fmt.Sprintf("TestAddInScalarMat3 %d %d", testIndex, k))
				break
			}
		}
	}
}

func TestSubScalarMat3(t *testing.T) {
	cases := []struct {
		orig  [9]float64
		value float64
	}{
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, 1},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, -1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, -1},
	}

	m := Mat3{}
	for testIndex, c := range cases {
		m.Load(c.orig)
		ret_mat := m.SubScalar(c.value)

		get := ret_mat.Dump()
		for k, _ := range c.orig {
			if get[k] != c.orig[k]-c.value {
				t.Error(fmt.Sprintf("TestSubScalarMat3 %d %d", testIndex, k))
				break
			}
		}

		ret_mat2 := m.SubInScalar(c.value)
		if ret_mat2 != &m {
			t.Error(fmt.Sprintf("TestSubInScalarMat3 %d", testIndex))
		}

		get = ret_mat2.Dump()
		for k, _ := range c.orig {
			if get[k] != c.orig[k]-c.value {
				t.Error(fmt.Sprintf("TestSubInScalarMat3 %d %d", testIndex, k))
				break
			}
		}
	}
}

func TestMultScalarMat3(t *testing.T) {
	cases := []struct {
		orig  [9]float64
		value float64
	}{
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, 0},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, 1},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, -1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, -1},
	}

	m := Mat3{}
	for testIndex, c := range cases {
		m.Load(c.orig)
		ret_mat := m.MultScalar(c.value)

		get := ret_mat.Dump()
		for k, _ := range c.orig {
			if get[k] != c.orig[k]*c.value {
				t.Error(fmt.Sprintf("TestMultScalarMat3 %d %d", testIndex, k))
				break
			}
		}

		ret_mat2 := m.MultInScalar(c.value)
		if ret_mat2 != &m {
			t.Error(fmt.Sprintf("TestMultInScalarMat3 %d", testIndex))
		}

		get = ret_mat2.Dump()
		for k, _ := range c.orig {
			if get[k] != c.orig[k]*c.value {
				t.Error(fmt.Sprintf("TestMultInScalarMat3 %d %d", testIndex, k))
				break
			}
		}
	}
}

func TestDivScalarMat3(t *testing.T) {
	cases := []struct {
		orig  [9]float64
		value float64
	}{
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, 1},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, -1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, -1},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, -2},
	}

	m := Mat3{}
	for testIndex, c := range cases {
		m.Load(c.orig)
		ret_mat := m.DivScalar(c.value)

		get := ret_mat.Dump()
		for k, _ := range c.orig {
			if get[k] != c.orig[k]/c.value {
				t.Error(fmt.Sprintf("TestDivScalarMat3 %d %d", testIndex, k))
				break
			}
		}

		ret_mat2 := m.DivInScalar(c.value)
		if ret_mat2 != &m {
			t.Error(fmt.Sprintf("TestDivInScalarMat3 %d", testIndex))
		}

		get = ret_mat2.Dump()
		for k, _ := range c.orig {
			if get[k] != c.orig[k]/c.value {
				t.Error(fmt.Sprintf("TestDivInScalarMat3 %d %d", testIndex, k))
				break
			}
		}
	}
}

func TestAddMat3(t *testing.T) {
	cases := []struct {
		orig, other [9]float64
	}{
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
			[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{-1, -2, -3, -4, -5, -6, -7, -8, -9}},
	}

	m := Mat3{}
	m2 := Mat3{}
	for testIndex, c := range cases {
		m.Load(c.orig)
		m2.Load(c.other)

		ret_mat := m.Add(m2)
		get := ret_mat.Dump()
		for k, _ := range c.orig {
			if closeEq(get[k], c.orig[k]+c.other[k], epsilon) == false {
				t.Error(fmt.Sprintf("TestAddMat3 %d %d", testIndex, k))
				break
			}
		}

		ret_mat2 := m.AddIn(m2)
		if ret_mat2 != &m {
			t.Error(fmt.Sprintf("TestAddInMat3 %d", testIndex))
		}

		get = ret_mat2.Dump()
		for k, _ := range c.orig {
			if closeEq(get[k], c.orig[k]+c.other[k], epsilon) == false {
				t.Error(fmt.Sprintf("TestAddInMat3 %d %d", testIndex, k))
				break
			}
		}
	}
}

func TestSubMat3(t *testing.T) {
	cases := []struct {
		orig, other [9]float64
	}{
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
			[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{-1, -2, -3, -4, -5, -6, -7, -8, -9}},
	}

	m := Mat3{}
	m2 := Mat3{}
	for testIndex, c := range cases {
		m.Load(c.orig)
		m2.Load(c.other)

		ret_mat := m.Sub(m2)
		get := ret_mat.Dump()
		for k, _ := range c.orig {
			if closeEq(get[k], c.orig[k]-c.other[k], epsilon) == false {
				t.Error(fmt.Sprintf("TestSubMat3 %d %d", testIndex, k))
				break
			}
		}

		ret_mat2 := m.SubIn(m2)
		if ret_mat2 != &m {
			t.Error(fmt.Sprintf("TestSubInMat3 %d", testIndex))
		}

		get = ret_mat2.Dump()
		for k, _ := range c.orig {
			if closeEq(get[k], c.orig[k]-c.other[k], epsilon) == false {
				t.Error(fmt.Sprintf("TestSubInMat3 %d %d", testIndex, k))
				break
			}
		}
	}
}

func TestMultMat3(t *testing.T) {
	cases := []struct {
		orig, other, want [9]float64
	}{
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}},

		{[9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},

		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}},

		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{30, 36, 42, 66, 81, 96, 102, 126, 150}},

		{[9]float64{5, 6, 7, 8, 1, 2, 3, 4, 13},
			[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{78, 96, 114, 26, 37, 48, 110, 130, 150}},

		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{5, 6, 7, 8, 1, 2, 3, 4, 13},
			[9]float64{30, 20, 50, 78, 53, 116, 126, 86, 182}},
	}

	orig := Mat3{}
	other := Mat3{}
	for testIndex, c := range cases {
		orig.Load(c.orig)
		other.Load(c.other)

		ret_mat := orig.Mult(other)
		get := ret_mat.Dump()
		for k, _ := range c.orig {
			if closeEq(get[k], c.want[k], epsilon) == false {
				t.Error(fmt.Sprintf("TestMultMat3 %d %d", testIndex, k))
				break
			}
		}

		ret_mat2 := orig.MultIn(other)
		if ret_mat2 != &orig {
			t.Error(fmt.Sprintf("TestMultInMat3 %d", testIndex))
		}
	}
}

func TestIdentityMat3(t *testing.T) {
	m := &Mat3{}
	m.Load([9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	m.ToIdentity()

	get := m.Dump()
	want := [9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}
	for k, _ := range get {
		if want[k] != get[k] {
			t.Error(fmt.Sprintf("TestIdentity %d", k))
			break
		}
	}
}

func TestTransposeMat3(t *testing.T) {
	cases := []struct {
		orig, want [9]float64
	}{
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
			[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}},

		{[9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
			[9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}},

		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			[9]float64{1, 4, 7, 2, 5, 8, 3, 6, 9}},
	}

	orig := &Mat3{}
	for testIndex, c := range cases {
		orig.Load(c.orig)
		get := orig.Transpose().Dump()

		for k, _ := range c.want {
			if closeEq(get[k], c.want[k], epsilon) == false {
				t.Error(fmt.Sprintf("TestTransposeMat3 %d %d", testIndex, k))
				break
			}
		}

		orig.Load(c.orig)
		get = orig.TransposeIn().Dump()
		for k, _ := range c.want {
			if closeEq(get[k], c.want[k], epsilon) == false {
				t.Error(fmt.Sprintf("TestTransposeInMat3 %d %d", testIndex, k))
				break
			}
		}
	}
}

func TestDeterminantMat3(t *testing.T) {
	want := 0.554020973727016224
	m := &Mat3{}
	m.Load([9]float64{
		0.5 * 0.5,
		-0.5 * 0.866,
		0.866,

		0.866*0.866*0.866 + 0.5*0.5,
		-0.866*0.866*0.866 + 0.5*0.5,
		0.866 * 0.5,

		-0.5*0.866*0.5 + 0.866*0.866,
		0.5*0.866*0.866 + 0.866*0.5,
		-0.5 * 0.5})

	get := m.Determinant()
	if closeEq(get, want, epsilon) == false {
		t.Error(fmt.Sprintf("TestDeterminantMat3 %v", get))
	}
}

func TestAdjointMat3(t *testing.T) {
	m := &Mat3{}
	m.Load([9]float64{
		0.5 * 0.5,
		-0.5 * 0.866,
		0.866,

		0.866*0.866*0.866 + 0.5*0.5,
		-0.866*0.866*0.866 + 0.5*0.5,
		0.866 * 0.5,

		-0.5*0.866*0.5 + 0.866*0.866,
		0.5*0.866*0.866 + 0.866*0.5,
		-0.5 * 0.5})

	m2 := m.Adjoint()
	get := m2.Dump()
	want := [9]float64{
		-0.249989, 0.591459, 0.158445,
		0.455852, -0.524473, 0.670684,
		0.939841, -0.432981, 0.289602}

	for k, _ := range want {
		// NOTE: the check uses lower precision because
		// the adjoint values I got from wolframalpha
		// only went up to 6 places
		if closeEq(get[k], want[k], 0.0001) == false {
			t.Error(fmt.Sprintf("TestAdjointMat3 %d %v %v", k, get[k], want[k]))
			break
		}
	}
}

func TestInverseMat3(t *testing.T) {
	m := &Mat3{}
	cases := []struct {
		orig, want        [9]float64
		want_inverse_flag bool
	}{
		{[9]float64{
			0.5 * 0.5,
			-0.5 * 0.866,
			0.866,

			0.866*0.866*0.866 + 0.5*0.5,
			-0.866*0.866*0.866 + 0.5*0.5,
			0.866 * 0.5,

			-0.5*0.866*0.5 + 0.866*0.866,
			0.5*0.866*0.866 + 0.866*0.5,
			-0.5 * 0.5}, [9]float64{
			-0.451227, 1.06758, 0.285991,
			0.822806, -0.946666, 1.21058,
			1.6964, -0.781524, 0.522727}, true},
		{[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
			[9]float64{0, 0, 0, 0, 0, 0, 0, 0, 0}, false},
		{[9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
			[9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}, true},
		{[9]float64{
			math.Cos(math.Pi / 2), -math.Sin(math.Pi / 2), 0,
			math.Sin(math.Pi / 2), math.Cos(math.Pi / 2), 0,
			0, 0, 1},
			[9]float64{
				math.Cos(math.Pi / 2), math.Sin(math.Pi / 2), 0,
				-math.Sin(math.Pi / 2), math.Cos(math.Pi / 2), 0,
				0, 0, 1}, true},
	}

	for testIndex, c := range cases {
		m.Load(c.orig)
		get_inverse_flag := m.HasInverse()
		if get_inverse_flag != c.want_inverse_flag {
			t.Error(fmt.Sprintf("TestInverseMat3 %d %v", testIndex, get_inverse_flag))
			continue
		}
		if get_inverse_flag == false {
			continue
		}

		m2 := m.Inverse()
		get := m2.Dump()
		for k, _ := range c.want {
			if closeEq(get[k], c.want[k], 0.0001) == false {
				t.Error(fmt.Sprintf("TestInverseMat3 %d %d %v %v", testIndex, k, get[k], c.want[k]))
				break
			}
		}
	}
}

func TestMultVec3Mat3(t *testing.T) {
	cases := []struct {
		orig_mat     [9]float64
		orig_v, want Vec3
	}{
		{[9]float64{1, 0, 0, 0, 1, 0, 0, 0, 1}, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{[9]float64{2, 0, 0, 0, 2, 0, 0, 0, 2}, Vec3{1, 0, 0}, Vec3{2, 0, 0}},
		{[9]float64{2, 0, 0, 0, 2, 0, 0, 0, 2}, Vec3{1, 1, 1}, Vec3{2, 2, 2}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, Vec3{1, 0, 0}, Vec3{1, 4, 7}},
		{[9]float64{1, 2, 3, 4, 5, 6, 7, 8, 9}, Vec3{1, 2, 3}, Vec3{14, 32, 50}},
	}

	m := Mat3{}
	for testIndex, c := range cases {
		m.Load(c.orig_mat)
		get := m.MultVec3(c.orig_v)
		if get.Eq(c.want) == false {
			t.Error(fmt.Sprintf("TestMultVec3Mat3 %d \n%v\n%v\n\n", testIndex, m, get))
		}
	}
}

func TestFromAxisAngleMat3(t *testing.T) {
	cases := []struct {
		angle     float64
		axis      Vec3
		start_vec Vec3
		want      Vec3
	}{

		//test basic rotations using a [1,0,0] vector
		{90, Vec3{0, 1, 0}, Vec3{1, 0, 0}, Vec3{0, 0, -1}},
		{90, Vec3{1, 0, 0}, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{90, Vec3{0, 0, 1}, Vec3{1, 0, 0}, Vec3{0, 1, 0}},
		{-90, Vec3{0, 1, 0}, Vec3{1, 0, 0}, Vec3{0, 0, 1}},
		{-90, Vec3{1, 0, 0}, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{-90, Vec3{0, 0, 1}, Vec3{1, 0, 0}, Vec3{0, -1, 0}},
		{360, Vec3{0, 0, 1}, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{180, Vec3{0, 0, 1}, Vec3{1, 0, 0}, Vec3{-1, 0, 0}}, //7

		//test basic rotations using a [0,1,0] vector
		{90, Vec3{0, 1, 0}, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{90, Vec3{1, 0, 0}, Vec3{0, 1, 0}, Vec3{0, 0, 1}},
		{90, Vec3{0, 0, 1}, Vec3{0, 1, 0}, Vec3{-1, 0, 0}},
		{-90, Vec3{0, 1, 0}, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{-90, Vec3{1, 0, 0}, Vec3{0, 1, 0}, Vec3{0, 0, -1}},
		{-90, Vec3{0, 0, 1}, Vec3{0, 1, 0}, Vec3{1, 0, 0}},
		{360, Vec3{0, 0, 1}, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{180, Vec3{0, 0, 1}, Vec3{0, 1, 0}, Vec3{0, -1, 0}}, //15

		// test negative axes
		{90, Vec3{0, -1, 0}, Vec3{1, 0, 0}, Vec3{0, 0, 1}},
		{90, Vec3{-1, 0, 0}, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{90, Vec3{0, 0, -1}, Vec3{1, 0, 0}, Vec3{0, -1, 0}},
		{-90, Vec3{0, -1, 0}, Vec3{1, 0, 0}, Vec3{0, 0, -1}},
		{-90, Vec3{-1, 0, 0}, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{-90, Vec3{0, 0, -1}, Vec3{1, 0, 0}, Vec3{0, 1, 0}},
		{360, Vec3{0, 0, -1}, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{180, Vec3{0, 0, -1}, Vec3{1, 0, 0}, Vec3{-1, 0, 0}}, //23

		// test arbitraty axis
		{360, Vec3{1, 1, 0}, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{90, Vec3{1, 1, 0}, Vec3{1, 0, 0}, Vec3{0.5, 0.5, -0.7071067811}},
		{45, Vec3{1, 1, 0}, Vec3{1, 0, 0}, Vec3{0.85355339059, 0.1464466094067, -0.5}}, //26
	}

	m := &Mat3{}
	for testIndex, c := range cases {
		c.axis.NormalizeIn()
		m.FromAxisAngle(Radians(c.angle), c.axis.X, c.axis.Y, c.axis.Z)

		get := m.MultVec3(c.start_vec)
		if get.Eq(c.want) == false {
			t.Error(fmt.Sprintf("TestFromAxisAngleMat3 %d \n%v\n%v\n\n", testIndex, m, get))
		}
	}
}

func TestAxisAngleMat3(t *testing.T) {
	cases := []struct {
		angle, x, y, z float64
	}{
		//test basic rotations using a [1,0,0] vector
		{90, 1, 0, 0},
		{90, 0, 1, 0},
		{90, 0, 0, 1},
		{45, 1, 0, 0},
		{45, 0, 1, 0},
		{45, 0, 0, 1}, //5
		{180, 1, 0, 0},
		{180, 0, 1, 0},
		{180, 0, 0, 1},
		{90, 1, 1, 0},
		{90, 1, 1, 0}, //10
		{90, 0, -1, 1},
		{45, 1, 0, 1},
		{45, 0, 1, 0},
		{45, 1, 0, 1},
		{180, 1, -2, 0}, //15
		{180, 0, 1, 20},
		{180, 0, 20, 1},
		{180, -4, 4, 1},
	}

	m := &Mat3{}
	for testIndex, c := range cases {
		v := Vec3{c.x, c.y, c.z}
		v.NormalizeIn()
		m.FromAxisAngle(Radians(c.angle), v.X, v.Y, v.Z)
		get_angle, get_x, get_y, get_z := m.AxisAngle()

		if !closeEq(Degrees(get_angle), c.angle, epsilon) ||
			!closeEq(get_x, v.X, epsilon) ||
			!closeEq(get_y, v.Y, epsilon) ||
			!closeEq(get_z, v.Z, epsilon) {

			if closeEq(get_angle, math.Pi, epsilon) &&
				closeEq(math.Abs(get_x)-math.Abs(v.X), 0, epsilon) &&
				closeEq(math.Abs(get_y)-math.Abs(v.Y), 0, epsilon) &&
				closeEq(math.Abs(get_z)-math.Abs(v.Z), 0, epsilon) {
				continue
			} else {
				t.Error(fmt.Sprintf("TestAxisAngleMat3 %d %v \n%f %f %f %f\n%f %f %f %f\n",
					testIndex, v, Degrees(get_angle), get_x, get_y, get_z, c.angle, v.X, v.Y, v.Z))
			}
		}
	}
}

func TestFromEulerMat3(t *testing.T) {
	common_cases := []struct {
		pitch, yaw, roll float64
		start_vec        Vec3
		want             Vec3
	}{
		{180, 0, 0, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{0, 180, 0, Vec3{1, 0, 0}, Vec3{-1, 0, 0}},
		{0, 0, 180, Vec3{1, 0, 0}, Vec3{-1, 0, 0}}, //2
		{180, 0, 0, Vec3{0, 1, 0}, Vec3{0, -1, 0}},
		{0, 180, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{0, 0, 180, Vec3{0, 1, 0}, Vec3{0, -1, 0}}, //5
		{180, 0, 0, Vec3{0, 0, 1}, Vec3{0, 0, -1}},
		{0, 180, 0, Vec3{0, 0, 1}, Vec3{0, 0, -1}},
		{0, 0, 180, Vec3{0, 0, 1}, Vec3{0, 0, 1}}, //8

		{180, 0, 0, Vec3{-1, 0, 0}, Vec3{-1, 0, 0}},
		{0, 180, 0, Vec3{-1, 0, 0}, Vec3{1, 0, 0}},
		{0, 0, 180, Vec3{-1, 0, 0}, Vec3{1, 0, 0}}, //11
		{180, 0, 0, Vec3{0, -1, 0}, Vec3{0, 1, 0}},
		{0, 180, 0, Vec3{0, -1, 0}, Vec3{0, -1, 0}},
		{0, 0, 180, Vec3{0, -1, 0}, Vec3{0, 1, 0}}, //14
		{180, 0, 0, Vec3{0, 0, -1}, Vec3{0, 0, 1}},
		{0, 180, 0, Vec3{0, 0, -1}, Vec3{0, 0, 1}},
		{0, 0, 180, Vec3{0, 0, -1}, Vec3{0, 0, -1}}, //17

		{0, 0, 0, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{0, 0, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{0, 0, 0, Vec3{0, 0, 1}, Vec3{0, 0, 1}}, //2
		{45, 90, 90, Vec3{0, 0, 1}, Vec3{math.Sqrt(2) / 2, math.Sqrt(2) / 2, 0}},

		//test basic rotations using a [0,1,0] vector
		// pitch,yaw,roll
		{0, 0, 90, Vec3{0, 1, 0}, Vec3{-1, 0, 0}},
		{0, 90, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{90, 0, 0, Vec3{0, 1, 0}, Vec3{0, 0, 1}},
		{0, 0, -90, Vec3{0, 1, 0}, Vec3{1, 0, 0}},
		{0, -90, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{-90, 0, 0, Vec3{0, 1, 0}, Vec3{0, 0, -1}},
		{0, 180, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}}, //6

		// test basic rotation using a [1,0,0] vector
		{0, 0, 90, Vec3{1, 0, 0}, Vec3{0, 1, 0}},
		{0, 90, 0, Vec3{1, 0, 0}, Vec3{0, 0, -1}},
		{90, 0, 0, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{0, 0, -90, Vec3{1, 0, 0}, Vec3{0, -1, 0}},
		{0, -90, 0, Vec3{1, 0, 0}, Vec3{0, 0, 1}},
		{-90, 0, 0, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{0, 0, 180, Vec3{1, 0, 0}, Vec3{-1, 0, 0}}, //13

		// basic rotation using a non major axis vector
		{0, 0, 90, Vec3{1, 1, 0}, Vec3{-1, 1, 0}},
		{0, 90, 0, Vec3{1, -1, 0}, Vec3{0, -1, -1}},
		{90, 0, 0, Vec3{-1, -1, 0}, Vec3{-1, 0, -1}}, //16

		// two rotations
		{90, 0, 45, Vec3{0, 0, 1}, Vec3{math.Sqrt(2) / 2, -math.Sqrt(2) / 2, 0}},
		{90, 45, 0, Vec3{0, 0, 1}, Vec3{0, -1, 0}},
		{45, 90, 0, Vec3{0, 0, 1}, Vec3{math.Sqrt(2) / 2, -math.Sqrt(2) / 2, 0}},
		{45, 90, 90, Vec3{0, 0, 1}, Vec3{math.Sqrt(2) / 2, math.Sqrt(2) / 2, 0}}, //20
	}

	m := Mat3{}
	for testIndex, c := range common_cases {
		// m = EulerToMat3(Radians(c.yaw), Radians(c.pitch), Radians(c.roll))
		m.FromEuler(Radians(c.pitch), Radians(c.yaw), Radians(c.roll))
		get := m.MultVec3(c.start_vec)
		if get.Eq(c.want) == false {
			t.Error(fmt.Sprintf("TestFromEulerMat3 %d \n%v\n%v\n\n", testIndex, m, get))
		}
	}
}

func TestEulerMat3(t *testing.T) {
	common_cases := []struct {
		pitch, yaw, roll float64
		start_vec        Vec3
		want             Vec3
	}{
		{180, 0, 0, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{0, 180, 0, Vec3{1, 0, 0}, Vec3{-1, 0, 0}},
		{0, 0, 180, Vec3{1, 0, 0}, Vec3{-1, 0, 0}}, //2
		{180, 0, 0, Vec3{0, 1, 0}, Vec3{0, -1, 0}},
		{0, 180, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{0, 0, 180, Vec3{0, 1, 0}, Vec3{0, -1, 0}}, //5
		{180, 0, 0, Vec3{0, 0, 1}, Vec3{0, 0, -1}},
		{0, 180, 0, Vec3{0, 0, 1}, Vec3{0, 0, -1}},
		{0, 0, 180, Vec3{0, 0, 1}, Vec3{0, 0, 1}}, //8

		{180, 0, 0, Vec3{-1, 0, 0}, Vec3{-1, 0, 0}},
		{0, 180, 0, Vec3{-1, 0, 0}, Vec3{1, 0, 0}},
		{0, 0, 180, Vec3{-1, 0, 0}, Vec3{1, 0, 0}}, //11
		{180, 0, 0, Vec3{0, -1, 0}, Vec3{0, 1, 0}},
		{0, 180, 0, Vec3{0, -1, 0}, Vec3{0, -1, 0}},
		{0, 0, 180, Vec3{0, -1, 0}, Vec3{0, 1, 0}}, //14
		{180, 0, 0, Vec3{0, 0, -1}, Vec3{0, 0, 1}},
		{0, 180, 0, Vec3{0, 0, -1}, Vec3{0, 0, 1}},
		{0, 0, 180, Vec3{0, 0, -1}, Vec3{0, 0, -1}}, //17

		{0, 0, 0, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{0, 0, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{0, 0, 0, Vec3{0, 0, 1}, Vec3{0, 0, 1}}, //20
		{45, 90, 90, Vec3{0, 0, 1}, Vec3{math.Sqrt(2) / 2, math.Sqrt(2) / 2, 0}},

		//test basic rotations using a [0,1,0] vector
		// pitch,yaw,roll
		{0, 0, 90, Vec3{0, 1, 0}, Vec3{-1, 0, 0}}, //22
		{0, 90, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{90, 0, 0, Vec3{0, 1, 0}, Vec3{0, 0, 1}},
		{0, 0, -90, Vec3{0, 1, 0}, Vec3{1, 0, 0}},
		{0, -90, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}},
		{-90, 0, 0, Vec3{0, 1, 0}, Vec3{0, 0, -1}},
		{0, 180, 0, Vec3{0, 1, 0}, Vec3{0, 1, 0}}, //28

		// test basic rotation using a [1,0,0] vector
		{0, 0, 90, Vec3{1, 0, 0}, Vec3{0, 1, 0}}, //29
		{0, 90, 0, Vec3{1, 0, 0}, Vec3{0, 0, -1}},
		{90, 0, 0, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{0, 0, -90, Vec3{1, 0, 0}, Vec3{0, -1, 0}},
		{0, -90, 0, Vec3{1, 0, 0}, Vec3{0, 0, 1}},
		{-90, 0, 0, Vec3{1, 0, 0}, Vec3{1, 0, 0}},
		{0, 0, 180, Vec3{1, 0, 0}, Vec3{-1, 0, 0}}, //35

		// basic rotation using a non major axis vector
		{0, 0, 90, Vec3{1, 1, 0}, Vec3{-1, 1, 0}},
		{0, 90, 0, Vec3{1, -1, 0}, Vec3{0, -1, -1}},
		{90, 0, 0, Vec3{-1, -1, 0}, Vec3{-1, 0, -1}}, //38

		// two rotations
		{90, 0, 45, Vec3{0, 0, 1}, Vec3{math.Sqrt(2) / 2, -math.Sqrt(2) / 2, 0}}, //39
		{90, 45, 0, Vec3{0, 0, 1}, Vec3{0, -1, 0}},
		{45, 90, 0, Vec3{0, 0, 1}, Vec3{math.Sqrt(2) / 2, -math.Sqrt(2) / 2, 0}},
		{45, 90, 90, Vec3{0, 0, 1}, Vec3{math.Sqrt(2) / 2, math.Sqrt(2) / 2, 0}}, //42
	}

	m := Mat3{}
	for testIndex, c := range common_cases {
		m.FromEuler(Radians(c.pitch), Radians(c.yaw), Radians(c.roll))
		x, y, z := m.Euler()

		if closeEq(Degrees(x), c.pitch, epsilon) && closeEq(Degrees(y), c.yaw, epsilon) && closeEq(Degrees(z), c.roll, epsilon) {
			continue
		}

		// The euler angles we got back didn't match, but lets see if the rotation
		// matrix it makes is still equivalent
		m.FromEuler(x, y, z)
		get := m.MultVec3(c.start_vec)
		if get.Eq(c.want) {
			continue
		}

		t.Error(fmt.Sprintf("TestEulerMat3 %d %f %f %f", testIndex, x, y, z))
	}
}
