package objects

import (
	"gonum.org/v1/gonum/mat"
	"math"
	"raytracer/ray"
	"raytracer/util"
)

type Sphere struct {
	center *mat.VecDense
	radius float64
}

func NewSphere(center *mat.VecDense, radius float64) Sphere {
	return Sphere{
		center,
		radius,
	}
}

func (s *Sphere) Hit(ray *ray.Ray, interval util.Interval, hitRecord *ray.HitRecord) bool {
	oc := util.NewZeroVector()
	oc.SubVec(s.center, ray.Origin)

	a := mat.Dot(ray.Direction, ray.Direction)
	h := mat.Dot(ray.Direction, oc)
	c := mat.Dot(oc, oc) - s.radius*s.radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return false
	}
	sqrtd := math.Sqrt(discriminant)
	root := (h - sqrtd) / a

	if !interval.Surrounds(root) {
		root = (h + sqrtd) / a
		if !interval.Surrounds(root) {
			return false
		}
	}

	hitRecord.Time = root
	hitRecord.Point = ray.At(hitRecord.Time)

	normal := util.NewZeroVector()
	normal.SubVec(hitRecord.Point, s.center)
	util.MakeUnitVector(normal)
	hitRecord.SetFaceNormal(ray, normal)

	return true
}
