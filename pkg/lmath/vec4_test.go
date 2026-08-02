package lmath

import (
	"fmt"
	"math"
	"testing"
)

func TestEqualVec4(t *testing.T) {
	var cases = []struct {
		orig, other Vec4
		want        bool
	}{
		{Vec4{0, 0, 0, 1}, Vec4{1, 2, 3, 1}, false},
		{Vec4{1, 2, 3, 1}, Vec4{0, 0, 0, 1}, false},
		{Vec4{1, 2, 3, 1}, Vec4{-1, -2, -3, 1}, false},
		{Vec4{0, 0, 0, 1}, Vec4{0, 0, 0, 1}, true},
		{Vec4{1.0, 2.0, 3.0, 1.0}, Vec4{1.0, 2.0, 3.0, 1.0}, true},
	}

	for testIndex, test := range cases {
		get := test.orig.Eq(test.other)
		if get != test.want {
			t.Error(fmt.Sprintf("TestEqualVec4 %d", testIndex))
		}
	}
}

func TestAddVec4(t *testing.T) {

	var cases = []struct {
		orig, other, want Vec4
	}{
		{Vec4{0, 0, 0, 0}, Vec4{1, 2, 3, 4}, Vec4{1, 2, 3, 4}},
		{Vec4{1, 2, 3, 4}, Vec4{0, 0, 0, 0}, Vec4{1, 2, 3, 4}},
		{Vec4{1, 2, 3, 4}, Vec4{-1, -2, -3, -4}, Vec4{0, 0, 0, 0}},
		{Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 1}, Vec4{0, 0, 0, 1}},
	}

	for testIndex, test := range cases {
		get := test.orig.Add(test.other)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestAddVec4 %d %v", testIndex, get))
		}

		get2 := test.orig.AddIn(test.other)
		if get2 != &test.orig || get2.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestAddVec4 AddIn %d", testIndex))
		}
	}
}

func TestSubVec4(t *testing.T) {

	var cases = []struct {
		orig, other, want Vec4
	}{
		{Vec4{0, 0, 0, 0}, Vec4{1, 2, 3, 4}, Vec4{-1, -2, -3, -4}},
		{Vec4{1, 2, 3, 4}, Vec4{0, 0, 0, 0}, Vec4{1, 2, 3, 4}},
		{Vec4{1, 2, 3, 4}, Vec4{-1, -2, -3, -4}, Vec4{2, 4, 6, 8}},
		{Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 0}},
	}

	for testIndex, test := range cases {
		get := test.orig.Sub(test.other)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestSubVec4 %d", testIndex))
		}

		get2 := test.orig.SubIn(test.other)
		if get2 != &test.orig || get2.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestSubVec4 SubIn %d", testIndex))
		}
	}
}

func TestAddScalarVec4(t *testing.T) {
	cases := []struct {
		orig, want Vec4
		scale      float64
	}{
		{Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 0}, 0},
		{Vec4{0, 0, 0, 0}, Vec4{1, 1, 1, 1}, 1},
		{Vec4{0, 0, 0, 0}, Vec4{-1, -1, -1, -1}, -1},
		{Vec4{1, 2, 3, 4}, Vec4{5, 6, 7, 8}, 4},
		{Vec4{1, 2, 3, 4}, Vec4{-3, -2, -1, 0}, -4},
		{Vec4{1, -2, 3, -4}, Vec4{5, 2, 7, 0}, 4},
	}

	for testIndex, test := range cases {
		get := test.orig.AddScalar(test.scale)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestAddScalarVec4 %d", testIndex))
		}

		get2 := test.orig.AddInScalar(test.scale)
		if get2 != &test.orig || get2.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestAddInScalarVec4 %d", testIndex))
		}
	}
}

func TestSubScalarVec4(t *testing.T) {
	cases := []struct {
		orig, want Vec4
		scale      float64
	}{
		{Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 0}, 0},
		{Vec4{0, 0, 0, 0}, Vec4{-1, -1, -1, -1}, 1},
		{Vec4{0, 0, 0, 0}, Vec4{1, 1, 1, 1}, -1},
		{Vec4{1, 2, 3, 4}, Vec4{-3, -2, -1, 0}, 4},
		{Vec4{1, 2, 3, 4}, Vec4{5, 6, 7, 8}, -4}, //4
		{Vec4{1, -2, 3, -4}, Vec4{-3, -6, -1, -8}, 4},
	}

	for testIndex, test := range cases {
		get := test.orig.SubScalar(test.scale)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestSubScalarVec4 %d", testIndex))
		}

		get2 := test.orig.SubInScalar(test.scale)
		if get2 != &test.orig || get2.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestSubInScalarVec4 %d", testIndex))
		}
	}
}

func TestMultScalarVec4(t *testing.T) {
	var cases = []struct {
		orig  Vec4
		scale float64
		want  Vec4
	}{
		{Vec4{0, 0, 0, 0}, 2.0, Vec4{0, 0, 0, 0}},
		{Vec4{0, 0, 0, 0}, -2, Vec4{0, 0, 0, 0}},
		{Vec4{1, 2, 3, 4}, 2, Vec4{2, 4, 6, 8}},
		{Vec4{1, 2, 3, 4}, 0.5, Vec4{0.5, 1, 1.5, 2}},
		{Vec4{1, 2, 3, 4}, -1, Vec4{-1, -2, -3, -4}},
		{Vec4{1, 2, 3, 4}, -0.5, Vec4{-0.5, -1, -1.5, -2.0}},
		{Vec4{1, 2, 3, 4}, 0, Vec4{0, 0, 0, 0}},
	}

	for testIndex, test := range cases {
		get := test.orig.MultScalar(test.scale)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestMultVec4 %d", testIndex))
		}

		get2 := test.orig.MultInScalar(test.scale)
		if get2 != &test.orig || get2.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestMultVec4 MultIn %d", testIndex))
		}
	}
}

func TestDivScalarVec4(t *testing.T) {
	var cases = []struct {
		orig  Vec4
		scale float64
		want  Vec4
	}{
		{Vec4{0, 0, 0, 0}, 2.0, Vec4{0, 0, 0, 0}},
		{Vec4{0, 0, 0, 0}, -2.0, Vec4{0, 0, 0, 0}},
		{Vec4{1, 2, 3, 4}, 2, Vec4{0.5, 1, 1.5, 2.0}},
		{Vec4{1, 2, 3, 4}, 0.5, Vec4{2, 4, 6, 8}},
		{Vec4{1, 2, 3, 4}, -1, Vec4{-1, -2, -3, -4}},
		{Vec4{1, 2, 3, 4}, -0.5, Vec4{-2, -4, -6, -8}},
	}

	for testIndex, test := range cases {
		get := test.orig.DivScalar(test.scale)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestDivVec4 %d", testIndex))
		}

		get2 := test.orig.DivInScalar(test.scale)
		if get2 != &test.orig || get2.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestDivInVec4 %d", testIndex))
		}
	}
}

func TestOuterVec4(t *testing.T) {
	cases := []struct {
		orig, other, want Vec4
	}{
		{Vec4{0, 0, 0, 0}, Vec4{1, 2, 3, 4}, Vec4{0, 0, 0, 0}},
		{Vec4{1, 2, 3, 4}, Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 0}},
		{Vec4{1, 2, 3, 4}, Vec4{-1, -2, -3, -4}, Vec4{-1, -4, -9, -16}},
		{Vec4{1, 2, 3, 4}, Vec4{1, 2, 3, 4}, Vec4{1, 4, 9, 16}},
		{Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 0}},
	}

	for testIndex, test := range cases {
		get := test.orig.Outer(test.other)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestOuterVec4 %d", testIndex))
		}

		get2 := test.orig.OuterIn(test.other)
		if get2 != &test.orig || get2.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestOuterInVec4 %d", testIndex))
		}
	}
}

func TestDotVec4(t *testing.T) {
	var cases = []struct {
		orig  Vec4
		other Vec4
		want  float64
	}{
		{Vec4{0, 0, 0, 0}, Vec4{1, 2, 3, 4}, 0},
		{Vec4{1, 2, 3, 4}, Vec4{0, 0, 0, 0}, 0},
		{Vec4{1, 2, 3, 4}, Vec4{-1, -2, -3, -4}, -30},
		{Vec4{1, 2, 3, 4}, Vec4{1, 2, 3, 4}, 30},
		{Vec4{0, 0, 0, 0}, Vec4{0, 0, 0, 0}, 0},
	}

	for testIndex, test := range cases {
		get := test.orig.Dot(test.other)
		if get != test.want {
			t.Error(fmt.Sprintf("TestDotVec4 %d", testIndex))
		}
	}
}

func TestLengthVec4(t *testing.T) {
	var cases = []struct {
		orig Vec4
		want float64
	}{
		{Vec4{0, 0, 0, 0}, 0},
		{Vec4{1, 2, 3, 4}, math.Sqrt(30)},
		{Vec4{1, 0, 0, 0}, 1},
		{Vec4{1 / math.Sqrt(30), 2 / math.Sqrt(30), 3 / math.Sqrt(30), 4 / math.Sqrt(30)}, 1},
	}

	for testIndex, test := range cases {
		get := test.orig.Length()
		if !closeEq(get, test.want, epsilon) {
			t.Error(fmt.Sprintf("TestLengthVec4 %d %f", testIndex, get))
		}
	}
}

func TestNormalizeVec4(t *testing.T) {
	sqrt_30 := math.Sqrt(30)
	var cases = []struct {
		orig, want Vec4
	}{
		{Vec4{1, 2, 3, 4}, Vec4{1 / sqrt_30, 2 / sqrt_30, 3 / sqrt_30, 4 / sqrt_30}},
		{Vec4{1, 0, 0, 0}, Vec4{1, 0, 0, 0}},
		{Vec4{1 / sqrt_30, 2 / sqrt_30, 3 / sqrt_30, 4 / sqrt_30}, Vec4{1 / sqrt_30, 2 / sqrt_30, 3 / sqrt_30, 4 / sqrt_30}},
	}

	for testIndex, test := range cases {
		get := test.orig.Normalize()
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestNormalizeVec4 %d", testIndex))
		}

		get2 := test.orig.NormalizeIn()
		if get2 != &test.orig || get2.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestNormalizeInVec4 %d", testIndex))
		}
	}
}

func TestSetVec4(t *testing.T) {
	var cases = []struct {
		x, y, z, w float64
		orig, want Vec4
	}{
		{1, 2, 3, 4, Vec4{0, 0, 0, 0}, Vec4{1, 2, 3, 4}},
		{0, 0, 1, 0, Vec4{0, 0, 0, 0}, Vec4{0, 0, 1, 0}},
		{0, 0, 1, 0, Vec4{1, 2, 3, 4}, Vec4{0, 0, 1, 0}},
		{-1, 2, 4, 6, Vec4{1, -1, 3, 90}, Vec4{-1, 2, 4, 6}},
	}

	for testIndex, test := range cases {
		get := test.orig.Set(test.x, test.y, test.z, test.w)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestSetVec4 %d", testIndex))
		}
	}
}

func TestProjVec4(t *testing.T) {
	v := 0.808359542 / 2
	var cases = []struct {
		from, on, want Vec4
	}{
		{Vec4{1, 1, 0, 0}, Vec4{1, 0, 0, 0}, Vec4{math.Sqrt(2) / 2, 0, 0, 0}},
		{Vec4{1, 0, 0, 0}, Vec4{1, 0, 0, 0}, Vec4{1, 0, 0, 0}},
		// should probably be checkig [0,0,0] projected [1,0,0] => {NaN,NaN,NaN}
		{Vec4{20, 50, 3, 20}, Vec4{1, 1, 1, 1}, Vec4{v, v, v, v}},
	}

	for testIndex, test := range cases {
		get := test.from.Proj(test.on)
		if get.Eq(test.want) == false {
			t.Error(fmt.Sprintf("TestProjVec4 %+v %+v", testIndex, get))
		}
	}
}
