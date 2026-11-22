package vec

import (
	gomat "math"

	"github.com/StoneTrench/CovarexEngine/engine/core/math"
)

type Vec2[T math.Number] struct {
	X, Y T
}

func InitVec2FromXY[T math.Number](x, y T) Vec2[T] {
	return Vec2[T]{X: x, Y: y}
}
func InitVec2FromArray[T math.Number](slice [2]T) Vec2[T] {
	return Vec2[T]{
		X: slice[0],
		Y: slice[1],
	}
}
func InitVec2FromSlice[T math.Number](slice []T) Vec2[T] {
	return Vec2[T]{
		X: slice[0],
		Y: slice[1],
	}
}
func CastVec2[B, A math.Number](v Vec2[A]) Vec2[B] {
	return InitVec2FromXY(B(v.X), B(v.Y))
}

func (v Vec2[T]) SetX(x T) Vec2[T] {
	v.X = x
	return v
}
func (v Vec2[T]) SetY(y T) Vec2[T] {
	v.Y = y
	return v
}
func (v Vec2[T]) Set(x, y T) Vec2[T] {
	return Vec2[T]{
		X: x,
		Y: y,
	}
}

func (v Vec2[T]) ToArray() [2]T {
	return [2]T{v.X, v.Y}
}
func (v Vec2[T]) ToSlice() []T {
	return []T{v.X, v.Y}
}

func (v Vec2[T]) Scale(s T) Vec2[T] {
	v.X *= s
	v.Y *= s
	return v
}
func (v Vec2[T]) Divscale(s T) Vec2[T] {
	v.X /= s
	v.Y /= s
	return v
}
func (v Vec2[T]) Offset(x, y, z T) Vec2[T] {
	v.X += x
	v.Y += y
	return v
}

func (v Vec2[T]) Min(other Vec2[T]) Vec2[T] {
	v.X = min(v.X, other.X)
	v.Y = min(v.Y, other.Y)
	return v
}
func (v Vec2[T]) Max(other Vec2[T]) Vec2[T] {
	v.X = max(v.X, other.X)
	v.Y = max(v.Y, other.Y)
	return v
}
func (v Vec2[T]) Minmax(other ...Vec2[T]) (min, max Vec2[T]) {
	min, max = v, v
	for _, o := range other {
		min = min.Min(o)
		max = max.Max(o)
	}
	return min, max
}

func (v Vec2[T]) Add(other Vec2[T]) Vec2[T] {
	v.X += other.X
	v.Y += other.Y
	return v
}
func (v Vec2[T]) Sub(other Vec2[T]) Vec2[T] {
	v.X -= other.X
	v.Y -= other.Y
	return v
}

func (v Vec2[T]) Cross(other Vec2[T]) T {
	return v.X*other.Y - v.Y*other.X
}
func (v Vec2[T]) Dot(other Vec2[T]) T {
	return v.X*other.X + v.Y*other.Y
}
func (v Vec2[T]) Distance(other Vec2[T]) T {
	return v.Sub(other).Length()
}
func (v Vec2[T]) Manhattan(other Vec2[T]) T {
	return v.Sub(other).Abs().Sum()
}
func (v Vec2[T]) Normalize() Vec2[T] {
	return v.Divscale(v.Length())
}

func (v Vec2[T]) InnerProduct() T {
	return v.Dot(v)
}
func (v Vec2[T]) Length() T {
	return T(gomat.Sqrt(float64(v.InnerProduct())))
}

func (v Vec2[T]) Volume() T {
	return v.X * v.Y
}
func (v Vec2[T]) Sum() T {
	return v.X + v.Y
}

func (v Vec2[T]) Negate() Vec2[T] {
	v.X = -v.X
	v.Y = -v.Y
	return v
}
func (v Vec2[T]) Abs() Vec2[T] {
	v.X = math.Abs(v.X)
	v.Y = math.Abs(v.Y)
	return v
}
func (v Vec2[T]) Floor() Vec2[T] {
	v.X = T(gomat.Floor(float64(v.X)))
	v.Y = T(gomat.Floor(float64(v.Y)))
	return v
}
func (v Vec2[T]) Ceil() Vec2[T] {
	v.X = T(gomat.Ceil(float64(v.X)))
	v.Y = T(gomat.Ceil(float64(v.Y)))
	return v
}
func (v Vec2[T]) Round() Vec2[T] {
	v.X = T(gomat.Round(float64(v.X)))
	v.Y = T(gomat.Round(float64(v.Y)))
	return v
}

func (v Vec2[T]) Equals(other Vec2[T]) bool {
	return v.X == other.X && v.Y == other.Y
}
