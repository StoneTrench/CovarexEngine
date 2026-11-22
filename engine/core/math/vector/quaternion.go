package vec

import (
	gomat "math"

	"github.com/StoneTrench/CovarexEngine/engine/core/math"
)

// Quat is a unit quaternion on a 4d hypersphere.
//
// W is the real component.
// X Y and Z are the complex components i j and k.
//
// It is also just a Vec4, with most of the utility functions of a vector.
type Quat[T math.Number] struct {
	X, Y, Z, W T
}

func InitQuatIdentity[T math.Number]() Quat[T] {
	return Quat[T]{X: 0, Y: 0, Z: 0, W: 1}
}

// InitQuatAngleAxis creates a quaternion from an angle and an axis.
//
// x y z - have to be the elements of a unit vector.
func InitQuatAngleAxis[T math.Number](angle, x, y, z T) Quat[T] {
	sin64, cos64 := gomat.Sincos(float64(angle) / 2)
	sin, cos := T(sin64), T(cos64)
	return Quat[T]{
		X: x * sin,
		Y: y * sin,
		Z: z * sin,
		W: cos,
	}
}
func InitQuatWXYZ[T math.Number](x, y, z, w T) Quat[T] {
	return Quat[T]{X: x, Y: y, Z: z, W: w}
}

func CastQuat[B, A math.Number](a Quat[A]) Quat[B] {
	return Quat[B]{
		X: B(a.X),
		Y: B(a.Y),
		Z: B(a.Z),
		W: B(a.W),
	}
}

func (q Quat[T]) SetX(x T) Quat[T] {
	q.X = x
	return q
}
func (q Quat[T]) SetY(y T) Quat[T] {
	q.Y = y
	return q
}
func (q Quat[T]) SetZ(z T) Quat[T] {
	q.Z = z
	return q
}
func (q Quat[T]) SetW(w T) Quat[T] {
	q.W = w
	return q
}
func (q Quat[T]) Set(x, y, z, w T) Quat[T] {
	return Quat[T]{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}
}

func (q Quat[T]) ToArray() [4]T {
	return [4]T{q.X, q.Y, q.Z, q.W}
}
func (q Quat[T]) ToSlice() []T {
	return []T{q.X, q.Y, q.Z, q.W}
}

func (q Quat[T]) Scale(s T) Quat[T] {
	q.X *= s
	q.Y *= s
	q.Z *= s
	q.W *= s
	return q
}
func (q Quat[T]) Divscale(s T) Quat[T] {
	q.X /= s
	q.Y /= s
	q.Z /= s
	q.W /= s
	return q
}
func (q Quat[T]) Offset(x, y, z, w T) Quat[T] {
	q.X += x
	q.Y += y
	q.Z += z
	q.W += w
	return q
}

func (q Quat[T]) Min(other Quat[T]) Quat[T] {
	q.X = min(q.X, other.X)
	q.Y = min(q.Y, other.Y)
	q.Z = min(q.Z, other.Z)
	q.W = min(q.W, other.W)
	return q
}
func (q Quat[T]) Max(other Quat[T]) Quat[T] {
	q.X = max(q.X, other.X)
	q.Y = max(q.Y, other.Y)
	q.Z = max(q.Z, other.Z)
	q.W = max(q.W, other.W)
	return q
}
func (q Quat[T]) Minmax(other ...Quat[T]) (min, max Quat[T]) {
	min, max = q, q
	for _, o := range other {
		min = min.Min(o)
		max = max.Max(o)
	}
	return min, max
}

func (q Quat[T]) Add(other Quat[T]) Quat[T] {
	q.X += other.X
	q.Y += other.Y
	q.Z += other.Z
	q.W += other.W
	return q
}
func (q Quat[T]) Sub(other Quat[T]) Quat[T] {
	q.X -= other.X
	q.Y -= other.Y
	q.Z -= other.Z
	q.W -= other.W
	return q
}
func (q Quat[T]) Mul(o Quat[T]) Quat[T] {
	return Quat[T]{
		W: q.W*o.W - q.X*o.X - q.Y*o.Y - q.Z*o.Z,
		X: q.W*o.X - q.X*o.Y - q.Y*o.Z - q.Z*o.W,
		Y: q.W*o.Y - q.X*o.Z - q.Y*o.W - q.Z*o.X,
		Z: q.W*o.Z - q.X*o.W - q.Y*o.X - q.Z*o.Y,
	}
}

func (q Quat[T]) Dot(o Quat[T]) T {
	return q.X*o.X + q.Y*o.Y + q.Z*o.Z + q.W*o.W
}
func (q Quat[T]) Distance(other Quat[T]) T {
	return q.Sub(other).Length()
}
func (q Quat[T]) Manhattan(other Quat[T]) T {
	return q.Sub(other).Abs().Sum()
}
func (q Quat[T]) Normalize() Quat[T] {
	return q.Divscale(q.Length())
}

func (q Quat[T]) InnerProduct() T {
	return q.Dot(q)
}
func (q Quat[T]) Length() T {
	return T(gomat.Sqrt(float64(q.InnerProduct())))
}

func (q Quat[T]) Volume() T {
	return q.X * q.Y * q.Z * q.W
}
func (q Quat[T]) Sum() T {
	return q.X + q.Y + q.Z + q.W
}

func (q Quat[T]) Negate() Quat[T] {
	q.X = -q.X
	q.Y = -q.Y
	q.Z = -q.Z
	q.W = -q.W
	return q
}
func (q Quat[T]) Conjugate() Quat[T] {
	q.X = -q.X
	q.Y = -q.Y
	q.Z = -q.Z
	return q
}
func (q Quat[T]) Abs() Quat[T] {
	q.X = math.Abs(q.X)
	q.Y = math.Abs(q.Y)
	q.Z = math.Abs(q.Z)
	q.W = math.Abs(q.W)
	return q
}
func (q Quat[T]) Floor() Quat[T] {
	q.X = T(gomat.Floor(float64(q.X)))
	q.Y = T(gomat.Floor(float64(q.Y)))
	q.Z = T(gomat.Floor(float64(q.Z)))
	q.W = T(gomat.Floor(float64(q.W)))
	return q
}
func (q Quat[T]) Ceil() Quat[T] {
	q.X = T(gomat.Ceil(float64(q.X)))
	q.Y = T(gomat.Ceil(float64(q.Y)))
	q.Z = T(gomat.Ceil(float64(q.Z)))
	q.W = T(gomat.Ceil(float64(q.W)))
	return q
}
func (q Quat[T]) Round() Quat[T] {
	q.X = T(gomat.Round(float64(q.X)))
	q.Y = T(gomat.Round(float64(q.Y)))
	q.Z = T(gomat.Round(float64(q.Z)))
	q.W = T(gomat.Round(float64(q.W)))
	return q
}

func (q Quat[T]) Equals(other Quat[T]) bool {
	return q.X == other.X && q.Y == other.Y && q.Z == other.Z && q.W == other.W
}

func (q Quat[T]) ToMatrix() Mat4x4[T] {
	xx := q.X * q.X
	yy := q.Y * q.Y
	zz := q.Z * q.Z
	xy := q.X * q.Y
	xz := q.X * q.Z
	yz := q.Y * q.Z
	wx := q.W * q.X
	wy := q.W * q.Y
	wz := q.W * q.Z

	return Mat4x4[T]{
		1 - 2*(yy+zz),
		2 * (xy - wz),
		2 * (xz + wy),
		0,

		2 * (xy + wz),
		1 - 2*(xx+zz),
		2 * (yz - wx),
		0,

		2 * (xz - wy),
		2 * (yz + wx),
		1 - 2*(xx+yy),
		0,

		0,
		0,
		0,
		1,
	}
}

// MulVec3 is a passive rotation.
//
// Active rotation is when the point is rotated with respect to the coordinate system, and passive rotation is when the coordinate system is rotated with respect to the point. The two rotations are opposite from each other.
//
//	func (q Quaternion[T]) MulVec(v Vector3) Vector3 {
//		vQuat := Quaternion[T]{W: 0, X: v.X, Y: v.Y, Z: v.Z}
//		res := q.MulQuat(vQuat).MulQuat(q.Conjugate())
//		return Vector3{X: res.X, Y: res.Y, Z: res.Z}
//	}
func (q Quat[T]) MulVec3(v Vec3[T]) Vec3[T] {
	numX := q.X * 2
	numY := q.Y * 2
	numZ := q.Z * 2

	tx := numY*v.Z - numZ*v.Y
	ty := numZ*v.X - numX*v.Z
	tz := numX*v.Y - numY*v.X

	return Vec3[T]{
		X: v.X + q.W*tx + (q.Y*tz - q.Z*ty),
		Y: v.Y + q.W*ty + (q.Z*tx - q.X*tz),
		Z: v.Z + q.W*tz + (q.X*ty - q.Y*tx),
	}
}
