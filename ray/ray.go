package ray

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

type Ray struct {
	Origin    *mat.VecDense
	Direction *mat.VecDense
}

func (r *Ray) at(t float64) *mat.VecDense {
	scaled := mat.NewVecDense(3, nil)
	scaled.ScaleVec(t, r.Direction)
	scaled.AddVec(scaled, r.Origin)

	return scaled
}

func (r *Ray) Color() *mat.VecDense {
	sphereCenter := mat.NewVecDense(3, []float64{0, 0, -1.})
	t := hitSphere(sphereCenter, 0.5, r)
	if t > 0.0 {
		n := r.at(t)
		n.SubVec(n, sphereCenter)
		// Getting the unit vector hasn't been so complicated ever
		n.ScaleVec(1./n.Norm(2), n)

		color := mat.NewVecDense(3, []float64{n.AtVec(0) + 1, n.AtVec(1) + 1, n.AtVec(2) + 1})
		color.ScaleVec(0.5, color)

		return color
	}

	unitDir := mat.NewVecDense(3, nil)
	// Unit again
	unitDir.ScaleVec(1/r.Direction.Norm(2), r.Direction)
	a := 0.5 * (unitDir.AtVec(1) + 1.0)

	firstColor := mat.NewVecDense(3, []float64{1, 1, 1.})
	firstColor.ScaleVec(1.0-a, firstColor)

	secondColor := mat.NewVecDense(3, []float64{0.5, 0.7, 1.})
	secondColor.ScaleVec(a, firstColor)

	firstColor.AddVec(firstColor, secondColor)
	return firstColor
}

// func hitSphere(center **mat.VecDense, radius float64, r *Ray) float64 {
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

func hitSphere(center *mat.VecDense, radius float64, r *Ray) float64 {
	oc := mat.NewVecDense(3, nil)
	oc.SubVec(center, r.Origin)

	a := mat.Dot(r.Direction, r.Direction)
	h := mat.Dot(r.Direction, oc)
	c := mat.Dot(oc, oc) - radius*radius

	discriminant := h*h - a*c

	if discriminant < 0 {
		return -1
	} else {
		return h - math.Sqrt(discriminant)/a
	}
}
