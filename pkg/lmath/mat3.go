package lmath

import (
	"fmt"
	"math"
)

const (
	mat3Dim = 3
)

var (
	Mat3Identity = Mat3{[9]float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}}
)

type Mat3 struct {
	mat [9]float64
}

// New Mat3 with the given values.
// Row-Order.
func NewMat3(
	m11, m12, m13,
	m21, m22, m23,
	m31, m32, m33 float64) *Mat3 {

	// 0   1   2   3
	// 4   5   6   7
	// 8   9   10  11
	// 12  13  14  15
	out := Mat3{}
	out.mat[0] = m11
	out.mat[1] = m12
	out.mat[2] = m13
	out.mat[3] = m21
	out.mat[4] = m22
	out.mat[5] = m23
	out.mat[6] = m31
	out.mat[7] = m32
	out.mat[8] = m33
	return &out
}

// Load the matrix with 9 floats.
// Specified in Row-Major order.
func (this *Mat3) Load(m [9]float64) *Mat3 {
	this.mat = m
	return this
}

// Load the matrix with 9 floats.
// Specified in Row-Major order.
func (this *Mat3) Load32(m [9]float32) *Mat3 {
	for k, v := range m {
		this.mat[k] = float64(v)
	}
	return this
}

// Retrieve a 9 float array of all the values of the matrix.
// Returned in Row-Major order.
func (this Mat3) Dump() (m [9]float64) {
	m = this.mat
	return
}

// Retrieve a 9 float64 array of all the values of the matrix.
// Returned in Col-Major order.
func (this Mat3) DumpOpenGL() (m [9]float64) {
	m[0], m[1], m[2] = this.Col(0)
	m[3], m[4], m[5] = this.Col(1)
	m[6], m[7], m[8] = this.Col(2)
	return
}

// Retrieve a 9 float32 array of all the values of the matrix.
// Returned in Col-Major order.
func (this Mat3) DumpOpenGLf32() (m [9]float32) {
	m[0] = float32(this.mat[0])
	m[1] = float32(this.mat[3])
	m[2] = float32(this.mat[6])

	m[3] = float32(this.mat[1])
	m[4] = float32(this.mat[4])
	m[5] = float32(this.mat[7])

	m[6] = float32(this.mat[2])
	m[7] = float32(this.mat[5])
	m[8] = float32(this.mat[8])
	return
}

// Return a copy of this matrix.
// Carbon-copy of all elements
func (this Mat3) Copy() Mat3 {
	return this
}

// Compare this matrix to the other.
// Return true if all elements between them are the same.
// Equality is measured using an epsilon (< 0.0000001).
func (this Mat3) Eq(other Mat3) bool {
	for k, _ := range this.mat {
		if closeEq(this.mat[k], other.mat[k], epsilon) == false {
			return false
		}
	}
	return true
}

// Retrieve the element at row and column.
// 0 indexed.
// Does not do any bounds checking.
func (this Mat3) Get(row, col int) float64 {
	return this.mat[row*mat3Dim+col]
}

// Set the value at the specified column and row.
// 0 indexed.
// Does not do any bounds checking.
func (this *Mat3) Set(row, col int, value float64) *Mat3 {
	this.mat[row*mat3Dim+col] = value
	return this
}

// Retrieve the element at the given index assuming a linear array.
// (i.e matrix[0], matrix[5]).
// 0 indexed.
func (this Mat3) At(index int) float64 {
	return this.mat[index]
}

// Set the element of the matrix specified at the index to the given value.
// 0 indexed.
// Return a pointer to the 'this'
func (this *Mat3) SetAt(index int, value float64) *Mat3 {
	this.mat[index] = value
	return this
}

// Set the specified row of the matrix to the given x,y,z values.
// 0 indexed.
// Does not do bounds checking of the row.
func (this *Mat3) SetRow(row int, x, y, z float64) *Mat3 {
	this.mat[row*mat3Dim] = x
	this.mat[row*mat3Dim+1] = y
	this.mat[row*mat3Dim+2] = z
	return this
}

// Set the specified column of the matrix to the given x,y,z values.
// 0 indexed.
// Does not do bounds checking on the col.
func (this *Mat3) SetCol(col int, x, y, z float64) *Mat3 {
	this.mat[mat3Dim*0+col] = x
	this.mat[mat3Dim*1+col] = y
	this.mat[mat3Dim*2+col] = z
	return this
}

// Retrieve the x,y,z elements from the specified row.
// 0 indexed.
// Does not bounds check the row.
func (this Mat3) Row(row int) (x, y, z float64) {
	x = this.mat[row*mat3Dim]
	y = this.mat[row*mat3Dim+1]
	z = this.mat[row*mat3Dim+2]
	return
}

// Retrieve the x,y,z elements from the specified column.
// 0 indexed.
// Does not bounds check the column.
func (this Mat3) Col(col int) (x, y, z float64) {
	x = this.mat[mat3Dim*0+col]
	y = this.mat[mat3Dim*1+col]
	z = this.mat[mat3Dim*2+col]
	return
}

// Add in a constant value to all the terms fo the matrix.
// Return a new matrix with the result.
func (this Mat3) AddScalar(val float64) Mat3 {
	this.AddInScalar(val)
	return this
}

// Add in a constant value to all the terms fo the matrix.
// Returns a pointer to 'this'.
func (this *Mat3) AddInScalar(val float64) *Mat3 {
	for k, _ := range this.mat {
		this.mat[k] += val
	}
	return this
}

// Subtract in a constant value to all the terms fo the matrix.
// Return a new matrix with the result.
func (this Mat3) SubScalar(val float64) Mat3 {
	this.SubInScalar(val)
	return this
}

// Subtract in a constant value to all the terms fo the matrix.
// Returns a pointer to 'this'.
func (this *Mat3) SubInScalar(val float64) *Mat3 {
	for k, _ := range this.mat {
		this.mat[k] -= val
	}
	return this
}

// Multiplies in a constant value to all the terms fo the matrix.
// Return a new matrix with the result.
func (this Mat3) MultScalar(val float64) Mat3 {
	this.MultInScalar(val)
	return this
}

// Multiplies in a constant value to all the terms fo the matrix.
// Returns a pointer to 'this'.
func (this *Mat3) MultInScalar(val float64) *Mat3 {
	for k, _ := range this.mat {
		this.mat[k] *= val
	}
	return this
}

// Divides in a constant value to all the terms fo the matrix.
// Return a new matrix with the result.
//
//	precondition: val > 0
func (this Mat3) DivScalar(val float64) Mat3 {
	this.DivInScalar(val)
	return this
}

// Divides in a constant value to all the terms fo the matrix.
// Returns a pointer to 'this'.
//
//	precondition: val > 0
func (this *Mat3) DivInScalar(val float64) *Mat3 {
	for k, _ := range this.mat {
		this.mat[k] /= val
	}
	return this
}

// Adds the two matrices together ( ie.  this + other).
// Return a new matrix with the result.
func (this Mat3) Add(other Mat3) Mat3 {
	this.AddIn(other)
	return this
}

// Adds the two matrices together ( ie.  this + other).
// Stores the result in this.
// Returns this.
func (this *Mat3) AddIn(other Mat3) *Mat3 {
	for k, _ := range this.mat {
		this.mat[k] += other.mat[k]
	}
	return this
}

// Subtract the two matrices together ( ie.  this - other).
// Return a new matrix with the result.
func (this Mat3) Sub(other Mat3) Mat3 {
	this.SubIn(other)
	return this
}

// Subtract the two matrices together ( ie.  this - other).
// Stores the result in this.
// Returns this.
func (this *Mat3) SubIn(other Mat3) *Mat3 {
	for k, _ := range this.mat {
		this.mat[k] -= other.mat[k]
	}
	return this
}

// Multiply the two matrices together ( ie.  this * other).
// Return a new matrix with the result.
func (this Mat3) Mult(other Mat3) Mat3 {
	this.MultIn(other)
	return this
}

// Multiplies the two matrices together ( ie.  this * other).
// Stores the result in this.
// Returns this.
func (this *Mat3) MultIn(o Mat3) *Mat3 {
	// 0   1   2
	// 3   4   5
	// 6   7   8
	m := *this
	this.mat[0] = m.mat[0]*o.mat[0] + m.mat[1]*o.mat[3] + m.mat[2]*o.mat[6]
	this.mat[1] = m.mat[0]*o.mat[1] + m.mat[1]*o.mat[4] + m.mat[2]*o.mat[7]
	this.mat[2] = m.mat[0]*o.mat[2] + m.mat[1]*o.mat[5] + m.mat[2]*o.mat[8]

	this.mat[3] = m.mat[3]*o.mat[0] + m.mat[4]*o.mat[3] + m.mat[5]*o.mat[6]
	this.mat[4] = m.mat[3]*o.mat[1] + m.mat[4]*o.mat[4] + m.mat[5]*o.mat[7]
	this.mat[5] = m.mat[3]*o.mat[2] + m.mat[4]*o.mat[5] + m.mat[5]*o.mat[8]

	this.mat[6] = m.mat[6]*o.mat[0] + m.mat[7]*o.mat[3] + m.mat[8]*o.mat[6]
	this.mat[7] = m.mat[6]*o.mat[1] + m.mat[7]*o.mat[4] + m.mat[8]*o.mat[7]
	this.mat[8] = m.mat[6]*o.mat[2] + m.mat[7]*o.mat[5] + m.mat[8]*o.mat[8]

	return this
}

// Returns a new matrix which is transpose to this.
func (this Mat3) Transpose() Mat3 {
	this.TransposeIn()
	return this
}

// Take the transpose of this matrix.
func (this *Mat3) TransposeIn() *Mat3 {
	// TODO: can definitely be way more efficient
	// by only exchanging the column entries.
	m00, m01, m02 := this.Row(0)
	m10, m11, m12 := this.Row(1)
	m20, m21, m22 := this.Row(2)

	this.SetCol(0, m00, m01, m02)
	this.SetCol(1, m10, m11, m12)
	this.SetCol(2, m20, m21, m22)

	return this
}

// Get the determinant of the matrix.
// Uses a straight-up Cramers-Rule implementation.
func (this Mat3) Determinant() float64 {
	// 0   1   2
	// 3   4   5
	// 6   7   8

	// Use Cramer's rule to calculate the determinant
	return det3x3(
		this.mat[0], this.mat[1], this.mat[2],
		this.mat[3], this.mat[4], this.mat[5],
		this.mat[6], this.mat[7], this.mat[8])
}

// Returns a new matrix which is the Adjoint matrix of this.
func (this Mat3) Adjoint() Mat3 {
	a1, a2, a3 := this.mat[0], this.mat[1], this.mat[2]
	b1, b2, b3 := this.mat[3], this.mat[4], this.mat[5]
	c1, c2, c3 := this.mat[6], this.mat[7], this.mat[8]

	// 0 1 2 3          a1 a2 a3 a4
	// 4 5 6 7          b1 b2 b3 b4
	// 8 9 10 11        c1 c2 c3 c4
	// 12 13 14 15      d1 d2 d3 d4

	// 0 1 2          a1 a2 a3
	// 3 4 5          b1 b2 b3
	// 6 7 8          c1 c2 c3
	this.mat[0] = det2x2(b2, b3, c2, c3)
	this.mat[1] = -det2x2(b1, b3, c1, c3)
	this.mat[2] = det2x2(b1, b2, c1, c2)

	this.mat[3] = -det2x2(a2, a3, c2, c3)
	this.mat[4] = det2x2(a1, a3, c1, c3)
	this.mat[5] = -det2x2(a1, a2, c1, c2)

	this.mat[6] = det2x2(a2, a3, b2, b3)
	this.mat[7] = -det2x2(a1, a3, b1, b3)
	this.mat[8] = det2x2(a1, a2, b1, b2)

	this.TransposeIn()
	return this
}

// Returns a new matrix which is the inverse matrix of this.
// The bool flag is false if an inverse does not exist.
func (this Mat3) Inverse() Mat3 {
	// TODO: Needs further testing
	// Try out with rotation matrices.
	// The inverse of a valid rotation matrix should just be the transpose
	det := this.Determinant()
	return this.Adjoint().DivScalar(det)
}

// Returns true if the inverse of this matrix exists false otherwise.
// Internally it checks to see if the determinant is zero.
func (this Mat3) HasInverse() bool {
	return !closeEq(this.Determinant(), 0, epsilon)
}

// Sets the matrix to the identity matrix.
func (this *Mat3) ToIdentity() *Mat3 {
	this.mat = [9]float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
	return this
}

//==============================================================================

// Return true if the matrix is the identity matrix.
func (this Mat3) IsIdentity() bool {
	iden := [9]float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	}
	for k, _ := range iden {
		if !closeEq(this.mat[k], iden[k], 0) {
			return false
		}
	}
	return true
}

// Check to see if the matrix is a valid rotation matrix.
//
//	The two properties it checks are
//	1) Determinant() == 1
//	2) m*m.Transpose  == Identity
func (this Mat3) IsRotation() bool {
	return closeEq(this.Determinant(), 1, epsilon) && this.Mult(this.Transpose()).IsIdentity()
}

// Implement the Stringer interface
// Prints out each row of the matrix on its own line
func (this Mat3) String() string {
	return fmt.Sprintf("%f %f %f\n%f %f %f\n%f %f %f",
		this.mat[0], this.mat[1], this.mat[2],
		this.mat[3], this.mat[4], this.mat[5],
		this.mat[6], this.mat[7], this.mat[8])
}

// =============================================================================
// Create a 2D translation matrix for Mat3. Overwrites all values in the matrix.
func (this *Mat3) ToTranslate(x, y float64) *Mat3 {
	this.ToIdentity()
	this.Set(0, 2, x)
	this.Set(1, 2, y)
	return this
}

// Create a 2D scaling matrix for Mat3. Overwrites all values in the matrix.
func (this *Mat3) ToScale(x, y float64) *Mat3 {
	this.ToIdentity()
	this.Set(0, 0, x)
	this.Set(1, 1, y)
	return this
}

// Create a 2D shearing matrix for Mat3. Overwrites all values in the matrix.
//
//	0 x 0
//	y 0 0
//	0 0 1
func (this *Mat3) ToShear(x, y float64) *Mat3 {
	this.ToIdentity()
	this.Set(0, 0, 0)
	this.Set(1, 1, 0)

	this.Set(0, 1, x)
	this.Set(1, 0, y)
	return this
}

// Create a 2D rotation matrix about the Z axis
// cos  -sin   0
// sin   cos     0
// 0     0     1
func (this *Mat3) ToRotateZ(angle float64) *Mat3 {
	this.ToIdentity()
	this.Set(0, 0, math.Cos(angle))
	this.Set(0, 1, -math.Sin(angle))
	this.Set(1, 0, math.Sin(angle))
	this.Set(1, 1, math.Cos(angle))
	return this
}

//==============================================================================

// Multiplies the Vec3 against the matrix ( ie. result = Matrix * Vec).
// Returns a new vector with the result.
func (this Mat3) MultVec3(v Vec3) (out Vec3) {
	// 0   1   2
	// 3   4   5
	// 6   7   8
	out.Set(
		this.mat[0]*v.X+this.mat[1]*v.Y+this.mat[2]*v.Z,
		this.mat[3]*v.X+this.mat[4]*v.Y+this.mat[5]*v.Z,
		this.mat[6]*v.X+this.mat[7]*v.Y+this.mat[8]*v.Z,
	)
	return
}

// =============================================================================

// Return a rotation matrix which rotates a vector about the axis [x,y,z] with
// the given angle (radians).
// Set this matrix as a rotation matrix from the give angle(radians) and axis.
func (this *Mat3) FromAxisAngle(angle, x, y, z float64) *Mat3 {
	//Reference http://en.wikipedia.org/wiki/Rotation_matrix
	c := math.Cos(angle)
	s := math.Sin(angle)
	t := (1 - c)

	return this.Load([9]float64{c + x*x*t, x*y*t - z*s, x*z*t + y*s,
		y*x*t + z*s, c + y*y*t, y*z*t - x*s,
		z*x*t - y*s, z*y*t + x*s, c + z*z*t})
}

// Set this as a rotation matrix using the specified pitch,yaw, and roll paramters.
// Angles are in radians.
func (this *Mat3) FromEuler(pitch, yaw, roll float64) *Mat3 {
	cx := math.Cos(pitch)
	sx := math.Sin(pitch)
	cy := math.Cos(yaw)
	sy := math.Sin(yaw)
	cz := math.Cos(roll)
	sz := math.Sin(roll)

	// This matrix was derived by multiplying each indiviudual rotation matrix
	// together into a single matrix.
	// note the matrices are applied in reverse order compared to the application
	// of the rotations.
	//   roll         yaw             pitch
	// | cz  -sz  0 | | cy   0   sy | | 1    0    0  |
	// | sz   cz  0 |x| 0    1   0  |x| 0    cx  -sx |
	// | 0    0   1 | | -sy  0   cy | | 0    sx   cx |

	// first row
	this.mat[0] = cz * cy
	this.mat[1] = cz*sy*sx - sz*cx
	this.mat[2] = sz*sx + cz*cx*sy
	// second row
	this.mat[3] = sz * cy
	this.mat[4] = cz*cx + sx*sy*sz
	this.mat[5] = sz*sy*cx - cz*sx
	// third row
	this.mat[6] = -sy
	this.mat[7] = sx * cy
	this.mat[8] = cy * cx
	return this
}

// Return the axis (radians) and axis of this rotation matrix.
// Assumes the matrix is a valid rotation matrix.
func (this Mat3) AxisAngle() (angle, x, y, z float64) {
	// Reference
	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToAngle/
	m00, m01, m02 := this.Get(0, 0), this.Get(0, 1), this.Get(0, 2)
	m10, m11, m12 := this.Get(1, 0), this.Get(1, 1), this.Get(1, 2)
	m20, m21, m22 := this.Get(2, 0), this.Get(2, 1), this.Get(2, 2)

	if closeEq(math.Abs(m01-m10), 0, epsilon) &&
		closeEq(math.Abs(m02-m20), 0, epsilon) &&
		closeEq(math.Abs(m12-m21), 0, epsilon) {
		// singularity check
		// Checking for cases in which the angle is either 0 or 180

		if this.IsIdentity() {
			// If the angle is 0, then the rotation matrix will be the identity matrix
			// A 0 angle means that there is an arbitrary axis.
			angle, x, y, z = 0, 1, 0, 0
			return
		}

		// Angle is 180, we need to find the axis it rotates around
		angle = math.Pi

		xx := (m00 + 1) / 2
		yy := (m11 + 1) / 2
		zz := (m22 + 1) / 2
		xy := (m01 + m10) / 4
		xz := (m02 + m20) / 4
		yz := (m12 + m21) / 4

		if (xx > yy) && (xx > zz) { // m[0][0] is the largest diagonal term
			if xx < epsilon {
				x = 0
				y = math.Sqrt(2) / 2
				z = math.Sqrt(2) / 2
			} else {
				x = math.Sqrt(xx)
				y = xy / x
				z = xz / x
			}
		} else if yy > zz { // m[1][1] is the largest diagonal term
			if yy < epsilon {
				x = math.Sqrt(2) / 2
				y = 0
				z = math.Sqrt(2) / 2
			} else {
				y = math.Sqrt(yy)
				x = xy / y
				z = yz / y
			}
		} else { // m[2][2] is the largest diagonal term so base result on this
			if zz < epsilon {
				x = math.Sqrt(2) / 2
				y = math.Sqrt(2) / 2
				z = 0
			} else {
				z = math.Sqrt(zz)
				x = xz / z
				y = yz / z
			}
		}
		return
	}

	// no singularity; therefore calculate as normal
	angle = math.Acos((m00 + m11 + m22 - 1) / 2)
	A := (m21 - m12)
	B := (m02 - m20)
	C := (m10 - m01)

	x = A / math.Sqrt(A*A+B*B+C*C)
	y = B / math.Sqrt(A*A+B*B+C*C)
	z = C / math.Sqrt(A*A+B*B+C*C)
	return
}

// Return the pitch,yaw and roll values for the given rotation matrix.
// The returned euler angle may not be the exact angle in which you supplied
// but they can be used to make an equilvalent rotation matrix.
func (this Mat3) Euler() (pitch, yaw, roll float64) {
	// The method for calculating the euler angles from a rotation matrix
	// uses the method described in this document
	// http://staff.city.ac.uk/~sbbh653/publications/euler.pdf

	// The rotation matrix we are using will be of the following form
	// cos(x) is abbreviated as cx ( similarily sin(x) = sx)
	// This corresponds to the pitch => yaw => roll rotation matrix
	// cz*cy       cz*sy*sx - sz*cx         sz*sx + cz*cx*sy   | r11 r12 r13
	// sz*cy       cz*cx + sx*sy*sz         sz*sy*cx - cz*sx   | r21 r22 r23
	// -sy         sx*cy                    cx*cy              | r31 r32 r33

	// We want to determine the x,y,z angles
	// 1) Find the 'y' angle
	//      This is easily accomplished because term r31 is simply '-sin(y)'
	// 2) There are two possible angles for y because
	//      sin(y) == sin(pi - y)
	// 3) To find the value of x, we observe the following
	//      r32/r33 = tan(x)
	//      (sin(x)cos(y)) / (cos(x)cos(y))
	//      (sin(x)/cos(x)) == tan(x) by defn.
	// 4) We can also calculate x and z by.
	//      x = atan2(r32,r33) == atan2( (sin(x)cos(y)) / (cos(x)cos(y)) )
	//      z = atan2(r21,r11) == atan2( (sin(z)cos(y)) / (cos(z)cos(y)) )

	var x, y, z float64
	r31 := this.Get(2, 0)
	if closeEq(r31, 1, epsilon) {
		// we are in gimbal lock
		z = 0
		y = -math.Pi / 2
		x = -z + math.Atan2(-this.Get(0, 1), -this.Get(0, 2))
	} else if closeEq(r31, -1, epsilon) {
		// we are in gimbal lock
		z = 0
		y = math.Pi / 2
		x = z + math.Atan2(this.Get(0, 1), this.Get(0, 2))
	} else {
		y = -math.Asin(r31)
		cos_y := math.Cos(y)
		x = math.Atan2(this.Get(2, 1)/cos_y, this.Get(2, 2)/cos_y)
		z = math.Atan2(this.Get(1, 0)/cos_y, this.Get(0, 0)/cos_y)

		// There are two alternative values for y,here is the second option
		// y2 := math.Pi - y
		// cos_y2 := math.Cos(y2)
		// x2 := math.Atan2(this.Get(2, 1)/cos_y2, this.Get(2, 2)/cos_y2)
		// z2 := math.Atan2(this.Get(1, 0)/cos_y2, this.Get(0, 0)/cos_y2)
	}

	pitch = x
	yaw = y
	roll = z
	return
}

// Creates a rotation matrix from the given quaternion. Return this
func (this *Mat3) FromQuat(q *Quat) *Mat3 {
	*this = q.Mat3()
	return this
}

// Returns the quaternion represented by this rotation matrix.
func (this Mat3) Quat() Quat {
	q := Quat{}
	q.FromMat3(this)
	return q
}
