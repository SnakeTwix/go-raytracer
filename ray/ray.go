package ray

import (
	"gonum.org/v1/gonum/mat"
	"math"
	"raytracer/util"
)

type Ray struct {
	Origin    *mat.VecDense
	Direction *mat.VecDense
}

func (r *Ray) At(t float64) *mat.VecDense {
	scaled := mat.NewVecDense(3, nil)
	scaled.ScaleVec(t, r.Direction)
	scaled.AddVec(scaled, r.Origin)

	return scaled
}

func (r *Ray) Color(world Hittable, depth int) *mat.VecDense {
	if depth == 0 {
		// Black color
		return util.NewZeroVector()
	}

	hitRecord := NewHitRecord()

	if world.Hit(r, util.NewInterval(0.001, math.MaxFloat64), &hitRecord) {

		//color := mat.NewVecDense(3, []float64{hitRecord.Normal.AtVec(0) + 1, hitRecord.Normal.AtVec(1) + 1, hitRecord.Normal.AtVec(2) + 1})
		//color.ScaleVec(0.5, color)
		//return color

		direction := util.NewRandomUnitVector()
		direction.AddVec(direction, hitRecord.Normal)

		newRayBounce := Ray{
			Origin:    hitRecord.Point,
			Direction: direction,
		}

		color := newRayBounce.Color(world, depth-1)
		color.ScaleVec(0.5, color)
		return color
	}

	unitDir := util.NewZeroVector()
	unitDir.CopyVec(r.Direction)
	util.MakeUnitVector(unitDir)
	a := 0.5 * (unitDir.AtVec(1) + 1.0)

	firstColor := mat.NewVecDense(3, []float64{1, 1, 1.})
	firstColor.ScaleVec(1.0-a, firstColor)

	secondColor := mat.NewVecDense(3, []float64{0.5, 0.7, 1.})
	secondColor.ScaleVec(a, secondColor)

	firstColor.AddVec(firstColor, secondColor)
	return firstColor
}
