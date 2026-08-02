package lmath

import (
	"math"
)

// A Vector 3 containing the three components
// X, Y, Z
type Vec3 struct {
	X, Y, Z float64
}

var (
	Vec3Right   = Vec3{1, 0, 0}
	Vec3Up      = Vec3{0, 1, 0}
	Vec3Forward = Vec3{0, 0, 1}
	Vec3Zero    = Vec3{0, 0, 0}
)

// Returns a new vector which is the result of adding 'this' with the
// other vector
func (this Vec3) Add(other Vec3) Vec3 {
	this.AddIn(other)
	return this
}

// Adds 'this' with the other vector.
// Store the result into 'this'
// Return a pointer to 'this'
func (this *Vec3) AddIn(other Vec3) *Vec3 {
	this.X += other.X
	this.Y += other.Y
	this.Z += other.Z
	return this
}

// Returns a new vector which is the result of subtracting 'this' with the
// other vector
func (this Vec3) Sub(other Vec3) Vec3 {
	this.SubIn(other)
	return this
}

// Subtracts'this' with the other vector.
// Store the result into 'this'
// Return a pointer to 'this'
func (this *Vec3) SubIn(other Vec3) *Vec3 {
	this.X -= other.X
	this.Y -= other.Y
	this.Z -= other.Z
	return this
}

// Returns a new vector with the scalar added to every element
func (this Vec3) AddScalar(scale float64) Vec3 {
	this.AddInScalar(scale)
	return this
}

// Add the scale to every element in the vector
// Return this
func (this *Vec3) AddInScalar(scale float64) *Vec3 {
	this.X += scale
	this.Y += scale
	this.Z += scale
	return this
}

// Returns a new vector with the scalar subtracted to every element
func (this Vec3) SubScalar(scale float64) Vec3 {
	this.SubInScalar(scale)
	return this
}

// Subtract the scale from every element in the vector
// Return a pointer to 'this'
func (this *Vec3) SubInScalar(scale float64) *Vec3 {
	this.X -= scale
	this.Y -= scale
	this.Z -= scale
	return this
}

// Returns a new vector where every element is multiplied by the scale
func (this Vec3) MultScalar(scale float64) Vec3 {
	this.MultInScalar(scale)
	return this
}

// Multiply the each element of this vector with the scale value.
// Return a pointer to 'this'
func (this *Vec3) MultInScalar(scale float64) *Vec3 {
	this.X *= scale
	this.Y *= scale
	this.Z *= scale
	return this
}

// Returns a new vector where every element is division by the scale
func (this Vec3) DivScalar(scale float64) Vec3 {
	this.DivInScalar(scale)
	return this
}

// Divide the each element of this vector with the scale value.
// Return a pointer to 'this'
func (this *Vec3) DivInScalar(scale float64) *Vec3 {
	this.X /= scale
	this.Y /= scale
	this.Z /= scale
	return this
}

// Do a pair-wise element multiplication with the provided vector
// Returns a new vector with the result
func (this Vec3) Outer(other Vec3) Vec3 {
	this.OuterIn(other)
	return this
}

// Do a element-wise multiplication with the provided vector
// Store the result into 'this'
// Return a pointer to 'this'
func (this *Vec3) OuterIn(other Vec3) *Vec3 {
	this.X = this.X * other.X
	this.Y = this.Y * other.Y
	this.Z = this.Z * other.Z
	return this
}

// Returns the Dot product between 'this' and the other vector
func (this Vec3) Dot(other Vec3) float64 {
	return this.X*other.X +
		this.Y*other.Y +
		this.Z*other.Z
}

// Return the length of the vector
// sqrt(x^2 + y^2 + z^2)
func (this Vec3) Length() float64 {
	return math.Sqrt(this.X*this.X + this.Y*this.Y + this.Z*this.Z)
}

// Return the squared length of the vector
// x^2 + y^2 + z^2
func (this Vec3) LengthSq() float64 {
	return this.X*this.X + this.Y*this.Y + this.Z*this.Z
}

// Checks for equality between the vectors.
// Equal is all elements are equal within an epsilon ( < 0.0000001)
func (this Vec3) Eq(other Vec3) bool {
	return closeEq(this.X, other.X, epsilon) &&
		closeEq(this.Y, other.Y, epsilon) &&
		closeEq(this.Z, other.Z, epsilon)
}

// Checks for equality between the vectors.
// Equal is all elements are equal within an user specified e
func (this Vec3) CloseEq(other Vec3, e float64) bool {
	return closeEq(this.X, other.X, e) &&
		closeEq(this.Y, other.Y, e) &&
		closeEq(this.Z, other.Z, e)
}

// Return a new vector which is the normalized version of 'this'
func (this Vec3) Normalize() Vec3 {
	this.NormalizeIn()
	return this
}

// Normalize the vector
// Return a pointer to 'this'
func (this *Vec3) NormalizeIn() *Vec3 {
	mag := this.Length()
	return this.DivInScalar(mag)
}

// Set X,Y,Z parameters of the vector.
func (this *Vec3) Set(x, y, z float64) *Vec3 {
	this.X = x
	this.Y = y
	this.Z = z
	return this
}

// Make a vector which is the projection of this onto other
func (this Vec3) Proj(other Vec3) Vec3 {
	n := this.Length() * other.Length()
	return other.Normalize().MultScalar(this.Dot(other) / n)
}

// Return a copy of this vector
func (this Vec3) Copy() Vec3 {
	return this
}

// Retrieve all three x,y,z paramters at once
func (this Vec3) Dump() (float64, float64, float64) {
	return this.X, this.Y, this.Z
}

// Retrieve all three x,y,z paramters at once, returned as float32
func (this Vec3) Dumpf32() (float32, float32, float32) {
	return float32(this.X), float32(this.Y), float32(this.Z)
}

// convert to Vec4. Forth component is set to zero.
func (this Vec3) Vec4() Vec4 {
	return Vec4{this.X, this.Y, this.Z, 0}
}

//==============================================================================
// Vector 3 specific methods

// Returns a new vector which is the Cross product with 'this' X 'other'
func (this Vec3) Cross(other Vec3) Vec3 {
	this.CrossIn(other)
	return this
}

// Take the cross product between 'this' X 'other'
// Store the result into 'this'
// Return a pointer to 'this'
func (this *Vec3) CrossIn(other Vec3) *Vec3 {
	x := (this.Y*other.Z - other.Y*this.Z)
	y := -(this.X*other.Z - other.X*this.Z)
	z := (this.X*other.Y - other.X*this.Y)
	this.X, this.Y, this.Z = x, y, z
	return this
}

// Apply the matrix against the Vector
// Return a new vector with the result v*m
func (this Vec3) MultMat4(right Mat4) Vec3 {
	// 0   1   2   3
	// 4   5   6   7
	// 8   9   10  11
	// 12  13  14  15
	this.Set(
		this.X*right.At(0)+this.Y*right.At(4)+this.Z*right.At(8)+right.At(12),
		this.X*right.At(1)+this.Y*right.At(5)+this.Z*right.At(9)+right.At(13),
		this.X*right.At(2)+this.Y*right.At(6)+this.Z*right.At(10)+right.At(14),
	)
	return this
}
