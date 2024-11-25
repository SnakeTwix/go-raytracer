package ray

import (
	"math"
	"raytracer/vec3"
)

type Ray struct {
	Origin    *vec3.Point3
	Direction *vec3.Vec3
}

func (r *Ray) at(t float64) vec3.Point3 {
	scaled := r.Direction.MulNew(t)
	scaled.AddVec(r.Origin)

	return scaled
}

func (r *Ray) Color() vec3.Color {
	sphereCenter := vec3.Point3{Z: -1.0}
	t := hitSphere(&sphereCenter, 0.5, r)
	if t > 0.0 {
		n := r.at(t)
		n.SubVec(&sphereCenter)
		n.Unit()

		color := vec3.Color{X: n.X + 1, Y: n.Y + 1, Z: n.Z + 1}
		color.Mul(0.5)

		return color
	}

	unitDir := r.Direction.UnitNew()
	a := 0.5 * (unitDir.Y + 1.0)

	firstColor := vec3.Color{X: 1.0, Y: 1.0, Z: 1.0}
	firstColor.Mul(1.0 - a)

	secondColor := vec3.Color{X: 0.5, Y: 0.7, Z: 1.0}
	secondColor.Mul(a)

	firstColor.AddVec(&secondColor)

	return firstColor
}

// func hitSphere(center *vec3.Point3, radius float64, r *Ray) float64 {
// 	oc := center.SubVecNew(r.Origin)
// 	a := r.Direction.Dot(r.Direction)
// 	b := r.Direction.Dot(&oc) * -2.0
// 	c := oc.Dot(&oc) - radius*radius

// 	discriminant := b*b - 4*a*c

// 	if discriminant < 0 {
// 		return -1
// 	} else {
// 		return -b - math.Sqrt(discriminant)/(2.0*a)
// 	}
// }

func hitSphere(center *vec3.Point3, radius float64, r *Ray) float64 {
	oc := center.SubVecNew(r.Origin)
	a := r.Direction.LengthSquared()
	h := r.Direction.Dot(&oc)
	c := oc.LengthSquared() - radius*radius

	discriminant := h*h - a*c

	if discriminant < 0 {
		return -1
	} else {
		return h - math.Sqrt(discriminant)/a
	}
}
