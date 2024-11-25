package vec3

import (
	"fmt"
	"io"
	"log"
	"math"
)

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func (v *Vec3) Reverse() {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
}

func (v *Vec3) ReverseNew() Vec3 {
	return Vec3{
		X: -v.X,
		Y: -v.Y,
		Z: -v.Z,
	}
}

func (v *Vec3) AddVec(vec *Vec3) {
	v.X += vec.X
	v.Y += vec.Y
	v.Z += vec.Z
}

func (v *Vec3) AddVecNew(vec *Vec3) Vec3 {
	return Vec3{
		X: v.X + vec.X,
		Y: v.Y + vec.Y,
		Z: v.Z + vec.Z,
	}
}

func (v *Vec3) SubVec(vec *Vec3) {
	v.X -= vec.X
	v.Y -= vec.Y
	v.Z -= vec.Z
}

func (v *Vec3) SubVecNew(vec *Vec3) Vec3 {
	return Vec3{
		X: v.X - vec.X,
		Y: v.Y - vec.Y,
		Z: v.Z - vec.Z,
	}
}

func (v *Vec3) Mul(f float64) {
	v.X *= f
	v.Y *= f
	v.Z *= f
}

func (v *Vec3) MulNew(f float64) Vec3 {
	return Vec3{
		X: v.X * f,
		Y: v.Y * f,
		Z: v.Z * f,
	}
}

func (v *Vec3) MulVecNew(vec *Vec3) Vec3 {
	return Vec3{
		X: v.X * vec.X,
		Y: v.Y * vec.Y,
		Z: v.Z * vec.Z,
	}
}

func (v *Vec3) Div(f float64) {
	v.X /= f
	v.Y /= f
	v.Z /= f
}

func (v *Vec3) DivNew(f float64) Vec3 {
	return Vec3{
		X: v.X / f,
		Y: v.Y / f,
		Z: v.Z / f,
	}
}
func (v *Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v *Vec3) Dot(vec *Vec3) float64 {
	return v.X*vec.X + v.Y*vec.Y + v.Z*vec.Z
}

func (v *Vec3) Cross(vec *Vec3) Vec3 {
	return Vec3{
		X: v.Y*vec.Z - v.Z*vec.Y,
		Y: v.Z*vec.Y - v.Y*vec.Z,
		Z: v.X*vec.Y - v.Y*vec.X,
	}
}

func (v *Vec3) UnitNew() Vec3 {
	return v.DivNew(v.Length())
}

func (v *Vec3) Unit() {
	v.Div(v.Length())
}

func (v *Vec3) String() string {
	return fmt.Sprintf("%f %f %f", v.X, v.Y, v.Z)
}

type Point3 = Vec3
type Color = Vec3

func (v *Color) Write(buffer io.Writer) {
	r := v.X
	g := v.Y
	b := v.Z

	ir := int(255.999 * r)
	ig := int(255.999 * g)
	ib := int(255.999 * b)

	str := fmt.Sprintf("%d %d %d\n", ir, ig, ib)

	_, err := buffer.Write([]byte(str))

	if err != nil {
		log.Fatal("Failed to write to buffer: ", err)
	}

}
