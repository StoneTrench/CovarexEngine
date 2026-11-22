package vec

import (
	gomat "math"

	"github.com/StoneTrench/CovarexEngine/engine/core/math"
)

type Mat4x4[T math.Number] [16]T

func InitMat4x4Diagonal[T math.Number](x, y, z, w T) Mat4x4[T] {
	return Mat4x4[T]{
		x, 0, 0, 0,
		0, y, 0, 0,
		0, 0, z, 0,
		0, 0, 0, w,
	}
}
func InitMat4x4Identity[T math.Number]() Mat4x4[T] {
	return InitMat4x4Diagonal[T](1, 1, 1, 1)
}
func InitMat4x4Translation[T math.Number](x, y, z T) Mat4x4[T] {
	return Mat4x4[T]{
		1, 0, 0, x,
		0, 1, 0, y,
		0, 0, 1, z,
		0, 0, 0, 1,
	}
}
func InitMat4x4RotX[T math.Number](theta T) Mat4x4[T] {
	s64, c64 := gomat.Sincos(float64(theta))
	s, c := T(s64), T(c64)

	return Mat4x4[T]{
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	}
}
func InitMat4x4RotY[T math.Number](theta T) Mat4x4[T] {
	s64, c64 := gomat.Sincos(float64(theta))
	s, c := T(s64), T(c64)

	return Mat4x4[T]{
		c, 0, s, 0,
		0, 1, 0, 0,
		-s, 0, c, 0,
		0, 0, 0, 1,
	}
}
func InitMat4x4RotZ[T math.Number](theta T) Mat4x4[T] {
	s64, c64 := gomat.Sincos(float64(theta))
	s, c := T(s64), T(c64)

	return Mat4x4[T]{
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}
}
func InitMat4x4Perspective[T math.Number](fovy, width, height, near, far float64) Mat4x4[T] {
	tanHalfFovy := gomat.Tan(fovy / 2)
	aspect := width / height

	m00 := 1 / (aspect * tanHalfFovy)
	m11 := -(1 / tanHalfFovy)

	m22 := far / (near - far)
	m23 := -T(1)
	m32 := (far * near) / (near - far)

	return Mat4x4[T]{
		T(m00), 0, 0, 0,
		0, T(m11), 0, 0,
		0, 0, T(m22), T(m23),
		0, 0, T(m32), 0,
	}
}
func InitMat4x4LookAt[T math.Number](eye, center, up Vec3[T]) Mat4x4[T] {
	f := center.Sub(eye).Normalize()
	s := f.Cross(up.Normalize()).Normalize()
	u := s.Cross(f)

	return Mat4x4[T]{
		s.X, u.X, -f.X, 0,
		s.Y, u.Y, -f.Y, 0,
		s.Z, u.Z, -f.Z, 0,
		-s.Dot(eye), -u.Dot(eye), f.Dot(eye), 1,
	}
}
func CastMat4x4[B, A math.Number](v Mat4x4[A]) Mat4x4[B] {
	return Mat4x4[B]{
		B(v[0]), B(v[1]), B(v[2]), B(v[3]),
		B(v[4]), B(v[5]), B(v[6]), B(v[7]),
		B(v[8]), B(v[9]), B(v[10]), B(v[11]),
		B(v[12]), B(v[13]), B(v[14]), B(v[15]),
	}
}

func (a *Mat4x4[T]) Write(row, col int, value T) {
	a[(row<<2)+col] = value
}
func (a Mat4x4[T]) Get(row, col int) T {
	return a[(row<<2)+col]
}

func (a Mat4x4[T]) ToArray() [16]T {
	return a
}
func (a Mat4x4[T]) ToSlice() []T {
	return a[:]
}

func (a Mat4x4[T]) Scale(s T) Mat4x4[T] {
	for row := range 4 {
		ridx := row << 2
		for col := range 4 {
			index := ridx + col
			a[index] *= s
		}
	}
	return a
}
func (a Mat4x4[T]) Divscale(s T) Mat4x4[T] {
	for row := range 4 {
		ridx := row << 2
		for col := range 4 {
			index := ridx + col
			a[index] /= s
		}
	}
	return a
}

func (a Mat4x4[T]) Mul(b Mat4x4[T]) Mat4x4[T] {
	var result Mat4x4[T]

	for row := range 4 {
		for col := range 4 {
			var sum T

			for k := range 4 {
				sum += a.Get(row, k) * b.Get(k, col)
			}

			result.Write(row, col, sum)
		}
	}

	return result
}
func (a Mat4x4[T]) Add(b Mat4x4[T]) Mat4x4[T] {
	for row := range 4 {
		ridx := row << 2
		for col := range 4 {
			index := ridx + col
			a[index] += b[index]
		}
	}

	return a
}

func (a Mat4x4[T]) Transpose() Mat4x4[T] {
	return Mat4x4[T]{
		a[0], a[4], a[8], a[12],
		a[1], a[5], a[9], a[13],
		a[2], a[6], a[10], a[14],
		a[3], a[7], a[11], a[15],
	}
}
func (a Mat4x4[T]) HTMInverse() Mat4x4[T] {
	m11, m12, m13, m14 := a[0], a[1], a[2], a[3]
	m21, m22, m23, m24 := a[4], a[5], a[6], a[7]
	m31, m32, m33, m34 := a[8], a[9], a[10], a[11]

	a[3] = -(m11*m14 + m21*m24 + m31*m34)
	a[7] = -(m12*m14 + m22*m24 + m32*m34)
	a[11] = -(m13*m14 + m23*m24 + m33*m34)

	a[0], a[1], a[2] = m11, m21, m31
	a[4], a[5], a[6] = m12, m22, m32
	a[8], a[9], a[10] = m13, m23, m33

	a[12] = 0
	a[13] = 0
	a[14] = 0
	a[15] = 1

	return a
}

func (a Mat4x4[T]) MulVec3(v Vec3[T]) Vec3[T] {
	x, y, z := v.X, v.Y, v.Z
	v.X = x*a[0] + y*a[1] + z*a[2] + a[3]
	v.Y = x*a[4] + y*a[5] + z*a[6] + a[7]
	v.Y = x*a[8] + y*a[9] + z*a[10] + a[11]
	return v
}
