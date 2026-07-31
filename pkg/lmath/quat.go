package lmath

import (
	"math"
)

// References
//http://mathworld.wolfram.com/Quaternion.html
//https://luckytoilet.wordpress.com/2014/11/24/visualizing-quaternions-with-unity/#comments
//http://www.gamasutra.com/view/feature/131686/rotating_objects_using_quaternions.php?page=1

// Strutcure to hold the quaternion
// When using as a rotation quaternion, you must ensure that the quaternion is unit length.
// W is the angle in radians.
// X,Y,Z is the axis of rotation.
type Quat struct {
	W, X, Y, Z float64
}

var (
	QuatIdentity = Quat{1, 0, 0, 0}
	QuatZero     = Quat{0, 0, 0, 0}
)

// Return true if all elements are equal between both quaternions.
// Comparisons are done using an epsilon (< 0.0000001)
func (this Quat) Eq(other Quat) bool {
	return closeEq(this.X, other.X, epsilon) &&
		closeEq(this.Y, other.Y, epsilon) &&
		closeEq(this.Z, other.Z, epsilon) &&
		closeEq(this.W, other.W, epsilon)
}

// Set the component of the quaternion. Return this
func (this *Quat) Set(w, x, y, z float64) *Quat {
	this.W = w
	this.X = x
	this.Y = y
	this.Z = z
	return this
}

// Add the two quaternions (q + other).
// Returns a new quaternion with the result.
func (this Quat) Add(other Quat) Quat {
	this.AddIn(other)
	return this
}

// Add two quaternionns (q + other) storing the result in this. Returns this.
func (this *Quat) AddIn(other Quat) *Quat {
	this.W += other.W
	this.X += other.X
	this.Y += other.Y
	this.Z += other.Z
	return this
}

// Subtract the two quaternions (q - other).
// Returns a new quaternion with the result.
func (this Quat) Sub(other Quat) Quat {
	this.SubIn(other)
	return this
}

// Sub two quaternionns (q - other). Result is stored in this. Return this.
func (this *Quat) SubIn(other Quat) *Quat {
	this.W -= other.W
	this.X -= other.X
	this.Y -= other.Y
	this.Z -= other.Z
	return this
}

// Multiply the two quaterions.
// Return a new quaterion with the result.
// Multiplying quaternions is NOT commutative (order matters!).
func (this Quat) Mult(q Quat) Quat {
	this.MultIn(q)
	return this
}

// Multiply 'this' with the other quaternion.
// Store the result into 'this'.
// Return a pointer to 'this'.
// Multiplying quaternions is NOT commutative (order matters!).
func (this *Quat) MultIn(other Quat) *Quat {
	w := this.W*other.W - this.X*other.X - this.Y*other.Y - this.Z*other.Z
	x := this.W*other.X + this.X*other.W + this.Y*other.Z - this.Z*other.Y
	y := this.W*other.Y + this.Y*other.W - this.X*other.Z + this.Z*other.X
	z := this.W*other.Z + this.Z*other.W + this.X*other.Y - this.Y*other.X

	this.X = x
	this.Y = y
	this.Z = z
	this.W = w
	return this
}

// Add a scalar quantity to the quaternion.
// Returns a new quaternion with the result.
func (this Quat) AddScalar(val float64) Quat {
	this.AddInScalar(val)
	return this
}

// Add a scalar quantity to the quaternion. Returns this.
func (this *Quat) AddInScalar(val float64) *Quat {
	this.W += val
	this.X += val
	this.Y += val
	this.Z += val
	return this
}

// Subtract a scalar quantity to the quaternion.
// Returns a new quaternion with the result.
func (this Quat) SubScalar(val float64) Quat {
	this.SubInScalar(val)
	return this
}

// Add a scalar quantity to the quaternion. Returns this.
func (this *Quat) SubInScalar(val float64) *Quat {
	this.W -= val
	this.X -= val
	this.Y -= val
	this.Z -= val
	return this
}

// Multiply a scalar quantity to the quaternion.
// Returns a new quaternion with the result.
func (this Quat) MultScalar(val float64) Quat {
	this.MultInScalar(val)
	return this
}

// Add a scalar quantity to the quaternion. Returns this.
func (this *Quat) MultInScalar(val float64) *Quat {
	this.W *= val
	this.X *= val
	this.Y *= val
	this.Z *= val
	return this
}

// Divide a scalar quantity to the quaternion.
// Returns a new quaternion with the result.
func (this Quat) DivScalar(val float64) Quat {
	this.DivInScalar(val)
	return this
}

// Add a scalar quantity to the quaternion.
// Returns this.
func (this *Quat) DivInScalar(val float64) *Quat {
	this.W /= val
	this.X /= val
	this.Y /= val
	this.Z /= val
	return this
}

// Make this quaterion into a unit length quaternion.
// Returns a pointer to this.
func (this *Quat) ToUnit() *Quat {
	n := this.Norm()
	if closeEq(n, 0, epsilon) {
		return this
	}
	return this.DivInScalar(n)
}

// Returns the norm of this quaternion.
// sqrt(x^2 + y^2 + z^2 + w^2)
func (this Quat) Norm() float64 {
	return math.Sqrt(this.NormSq())
}

// Returns the squared norm of this quaternion.
// So that we can save a math.Sqrt.
// x^2 + y^2 + z^2 + w^2
func (this Quat) NormSq() float64 {
	return this.X*this.X +
		this.Y*this.Y +
		this.Z*this.Z +
		this.W*this.W
}

// Return a new quaternion which is the conjugate of this
// Conjugate is defined as [w,-x,-y,-z]
func (this Quat) Conjugate() Quat {
	this.ConjugateIn()
	return this
}

// Set the quaternion to the conjugate
// Conjugate is defined as [w,-x,-y,-z]
func (this *Quat) ConjugateIn() *Quat {
	this.X = -this.X
	this.Y = -this.Y
	this.Z = -this.Z
	return this
}

// Return a new quaternion which is the inverse of this
func (this Quat) Inverse() Quat {
	this.InverseIn()
	return this
}

// Set this quaternion as the inverse
func (this *Quat) InverseIn() *Quat {
	this.ConjugateIn()
	n := this.NormSq()
	if closeEq(n, 0, epsilon) {
		return this
	}
	return this.DivInScalar(n)
}

// Create a copy of this Quaternion and return the result.
func (this Quat) Copy() Quat {
	return this
}

//==============================================================================

// Return a new vector holding the axis component of the quaternion
func (this Quat) Axis() (out Vec3) {
	angle := 2 * math.Acos(this.W)
	s := math.Sin(angle / 2)
	out.X = this.X / s
	out.Y = this.Y / s
	out.Z = this.Z / s
	return
}

// return the angle (radians) component of the quaternion
func (this Quat) Angle() float64 {
	return 2 * math.Acos(this.W)
}

//==============================================================================

// Apply this quaternion as a rotation to the vec3.
// Perform the operation q * v * q^-1
// Return a new vector with result
func (this Quat) RotateVec3(v Vec3) (out Vec3) {
	vq := Quat{0.0, v.X, v.Y, v.Z}
	rs := this.Mult(vq).Mult(this.Inverse())
	out.Set(rs.X, rs.Y, rs.Z)
	return
}

// Create a quaternion from the specified euler angles.
// Return this
//
//	Perform the operation in the order.
//	pitch(x) => yaw(y) => row(z)
func (this *Quat) FromEuler(pitch, yaw, roll float64) *Quat {
	yawQ := Quat{math.Cos(yaw / 2), 0, math.Sin(yaw / 2), 0}
	pitchQ := Quat{math.Cos(pitch / 2), math.Sin(pitch / 2), 0, 0}
	rollQ := Quat{math.Cos(roll / 2), 0, 0, math.Sin(roll / 2)}

	// note must be applied in reverse order
	// pitch => yaw => roll
	*this = *(rollQ.MultIn(yawQ).MultIn(pitchQ))
	return this
}

// Set the quaternion as a rotation with with specified angle (radians) and axis.
// The axis should be normalized!
// Return this
func (this *Quat) FromAxisAngle(angle float64, x, y, z float64) *Quat {
	this.Set(math.Cos(angle/2),
		math.Sin(angle/2)*x,
		math.Sin(angle/2)*y,
		math.Sin(angle/2)*z)
	return this
}

// Extract out the euler angles from the quaternion
// Extract out the angles assuming the quaterion is encoded
// as pitch -> yaw -> roll
// The returned euler angle may not match the same angles passed in using the
// FromEuler(),but the angles are guaranteed to form an equivalent rotation quaternion.
func (this Quat) Euler() (pitch, yaw, roll float64) {
	// Reference http://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToEuler/

	test := this.X*this.Y + this.Z*this.W
	if test > 0.499 { // singularity at north pole
		yaw = 2 * math.Atan2(this.X, this.W)
		roll = math.Pi / 2
		pitch = 0
		return
	}
	if test < -0.499 { // singularity at south pole
		yaw = -2 * math.Atan2(this.X, this.W)
		roll = -math.Pi / 2
		pitch = 0
		return
	}

	sqx := this.X * this.X
	sqy := this.Y * this.Y
	sqz := this.Z * this.Z
	yaw = math.Atan2(2*this.Y*this.W-2*this.X*this.Z, 1-2*sqy-2*sqz)
	roll = math.Asin(2 * test)
	pitch = math.Atan2(2*this.X*this.W-2*this.Y*this.Z, 1-2*sqx-2*sqz)
	return
}

// Return the axis and angle of this rotation quaternion
func (this Quat) AxisAngle() (angle, x, y, z float64) {
	axis := this.Axis()
	return this.Angle(), axis.X, axis.Y, axis.Z
}

// Set the quaternion from the specified rotation matrix.
// Assumption is that the matrix is a valid rotation matrix.
// Matrix should be in right-hand coordinate system
// Pitch-Yaw-Roll euler angle formation
// Return this
func (this *Quat) fromMat(m [16]float64) *Quat {
	// Reference : http://www.flipcode.com/documents/matrfaq.html#Q55
	// 0  1  2  3
	// 4  5  6  7
	// 8  9  10 11
	// 12 13 14 15
	// trace := m.Get(0, 0) + m.Get(1, 1) + m.Get(2, 2) + 1
	trace := m[0] + m[5] + m[10] + 1

	if trace > 0 {
		s := 0.5 / math.Sqrt(trace)
		this.Set(
			0.25/s,
			(m[9]-m[6])*s,
			(m[2]-m[8])*s,
			(m[4]-m[1])*s,
		)
		return this
	}

	// Find the column which has the maximum diagonal value
	test_cols := [3]int{0, 5, 10}
	max_col := 0
	champ := m[test_cols[0]]
	for col := 1; col < 3; col += 1 {
		cand := m[test_cols[col]]
		//cand := m.Get(col, col)
		if cand > champ {
			champ = cand
			max_col = col
		}
	}

	// TODO : UNTESTED!!!
	var w, x, y, z, s float64
	switch max_col {
	case 0:
		s = 2 * math.Sqrt(1.0+m[0]-m[5]-m[10])
		x = 0.5 / 2
		y = (m[4] + m[1]) / s
		z = (m[8] + m[2]) / s
		w = (m[9] + m[6]) / s
	case 1:
		s = 2 * math.Sqrt(1.0+m[5]-m[0]-m[10])
		x = (m[4] + m[1]) / s
		y = 0.5 / 2
		z = (m[9] + m[6]) / s
		w = (m[8] + m[2]) / s
	case 2:
		s = 2 * math.Sqrt(1.0+m[10]-m[0]-m[5])
		x = (m[8] + m[2]) / s
		y = (m[9] + m[6]) / s
		z = 0.5 / 2
		w = (m[4] + m[1]) / s
	}

	this.Set(w, x, y, z)
	return this
}

// Return a mat4 from the provided quaternion
func (this Quat) mat() (m [16]float64) {
	// Reference
	// Derivation of the below matrix can be found here
	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/quaternionToMatrix/index.htm
	//     1 - 2y² - 2z²    2yx - 2wz        2xz + 2wy
	// M=  2xy + 2wz        1 - 2x² - 2z²    2yz - 2wx
	//     2xz - 2wy        2yz + 2wx        1 - 2x² - 2y²

	w, x, y, z := this.W, this.X, this.Y, this.Z

	// 0 1 2 3
	// 4 5 6 7
	// 8 9 10 11
	// 12 13 14 15
	m[0] = 1 - 2*y*y - 2*z*z
	m[1] = 2*y*z - 2*w*z
	m[2] = 2*x*z + 2*w*y
	m[3] = 0

	m[4] = 2*x*y + 2*w*z
	m[5] = 1 - 2*x*x - 2*z*z
	m[6] = 2*y*z - 2*w*x
	m[7] = 0

	m[8] = 2*x*z - 2*w*y
	m[9] = 2*y*z + 2*w*x
	m[10] = 1 - 2*x*x - 2*y*y
	m[11] = 0

	m[12] = 0
	m[13] = 0
	m[14] = 0
	m[15] = 1
	return m
}

func (this *Quat) FromMat4(m Mat4) *Quat {
	values := m.Dump()
	return this.fromMat(values)
}

func (this *Quat) FromMat3(m Mat3) *Quat {
	values := [16]float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	mxs := m.Dump()
	values[0] = mxs[0]
	values[1] = mxs[1]
	values[2] = mxs[2]
	values[4] = mxs[3]
	values[5] = mxs[4]
	values[6] = mxs[5]
	values[8] = mxs[6]
	values[9] = mxs[7]
	values[10] = mxs[8]

	return this.fromMat(values)
}

func (this Quat) Mat4() Mat4 {
	values := this.mat()
	m := Mat4{}
	m.Load(values)
	return m
}
func (this Quat) Mat3() Mat3 {
	v := this.mat()
	m := Mat3{}
	m.Load([9]float64{
		v[0], v[1], v[2],
		v[4], v[5], v[6],
		v[8], v[9], v[10],
	})
	return m
}
