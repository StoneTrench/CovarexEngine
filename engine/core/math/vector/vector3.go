package vec

import (
	gomat "math"

	"github.com/StoneTrench/CovarexEngine/engine/core/math"
)

type Vec3[T math.Number] struct {
	X, Y, Z T
}

func InitVec3[T math.Number](x, y, z T) Vec3[T] {
	return Vec3[T]{X: x, Y: y, Z: z}
}
func InitVec3FromArray[T math.Number](slice [3]T) Vec3[T] {
	return Vec3[T]{
		X: slice[0],
		Y: slice[1],
		Z: slice[2],
	}
}
func InitVec3FromSlice[T math.Number](slice []T) Vec3[T] {
	return Vec3[T]{
		X: slice[0],
		Y: slice[1],
		Z: slice[2],
	}
}
func CastVec3[B, A math.Number](v Vec3[A]) Vec3[B] {
	return InitVec3(B(v.X), B(v.Y), B(v.Z))
}

func (v Vec3[T]) SetX(x T) Vec3[T] {
	v.X = x
	return v
}
func (v Vec3[T]) SetY(y T) Vec3[T] {
	v.Y = y
	return v
}
func (v Vec3[T]) SetZ(z T) Vec3[T] {
	v.Z = z
	return v
}
func (v Vec3[T]) Set(x, y, z T) Vec3[T] {
	return Vec3[T]{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v Vec3[T]) ToArray() [3]T {
	return [3]T{v.X, v.Y, v.Z}
}
func (v Vec3[T]) ToSlice() []T {
	return []T{v.X, v.Y, v.Z}
}

func (v Vec3[T]) Scale(s T) Vec3[T] {
	v.X *= s
	v.Y *= s
	v.Z *= s
	return v
}
func (v Vec3[T]) Divscale(s T) Vec3[T] {
	v.X /= s
	v.Y /= s
	v.Z /= s
	return v
}
func (v Vec3[T]) Offset(x, y, z T) Vec3[T] {
	v.X += x
	v.Y += y
	v.Z += z
	return v
}

func (v Vec3[T]) Min(other Vec3[T]) Vec3[T] {
	v.X = min(v.X, other.X)
	v.Y = min(v.Y, other.Y)
	v.Z = min(v.Z, other.Z)
	return v
}
func (v Vec3[T]) Max(other Vec3[T]) Vec3[T] {
	v.X = max(v.X, other.X)
	v.Y = max(v.Y, other.Y)
	v.Z = max(v.Z, other.Z)
	return v
}
func (v Vec3[T]) Minmax(other ...Vec3[T]) (min, max Vec3[T]) {
	min, max = v, v
	for _, o := range other {
		min = min.Min(o)
		max = max.Max(o)
	}
	return min, max
}

func (v Vec3[T]) Add(other Vec3[T]) Vec3[T] {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}
func (v Vec3[T]) Sub(other Vec3[T]) Vec3[T] {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

func (v Vec3[T]) Cross(other Vec3[T]) Vec3[T] {
	return Vec3[T]{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}
func (v Vec3[T]) Dot(other Vec3[T]) T {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}
func (v Vec3[T]) Distance(other Vec3[T]) T {
	return v.Sub(other).Length()
}
func (v Vec3[T]) Manhattan(other Vec3[T]) T {
	return v.Sub(other).Abs().Sum()
}
func (v Vec3[T]) Normalize() Vec3[T] {
	return v.Divscale(v.Length())
}

func (v Vec3[T]) InnerProduct() T {
	return v.Dot(v)
}
func (v Vec3[T]) Length() T {
	return T(gomat.Sqrt(float64(v.InnerProduct())))
}

func (v Vec3[T]) Volume() T {
	return v.X * v.Y * v.Z
}
func (v Vec3[T]) Sum() T {
	return v.X + v.Y + v.Z
}

func (v Vec3[T]) Negate() Vec3[T] {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}
func (v Vec3[T]) Abs() Vec3[T] {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	v.Z = math.Abs(v.Z)
	return v
}
func (v Vec3[T]) Floor() Vec3[T] {
	v.X = T(gomat.Floor(float64(v.X)))
	v.Y = T(gomat.Floor(float64(v.Y)))
	v.Z = T(gomat.Floor(float64(v.Z)))
	return v
}
func (v Vec3[T]) Ceil() Vec3[T] {
	v.X = T(gomat.Ceil(float64(v.X)))
	v.Y = T(gomat.Ceil(float64(v.Y)))
	v.Z = T(gomat.Ceil(float64(v.Z)))
	return v
}
func (v Vec3[T]) Round() Vec3[T] {
	v.X = T(gomat.Round(float64(v.X)))
	v.Y = T(gomat.Round(float64(v.Y)))
	v.Z = T(gomat.Round(float64(v.Z)))
	return v
}

func (v Vec3[T]) Equals(other Vec3[T]) bool {
	return v.X == other.X && v.Y == other.Y && v.Z == other.Z
}
