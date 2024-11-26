package objects

import (
	"gonum.org/v1/gonum/mat"
	"math"
	"raytracer/ray"
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

func (s *Sphere) Hit(ray *ray.Ray, rayTmin float64, rayTmax float64, hitRecord *ray.HitRecord) bool {
	oc := mat.NewVecDense(3, nil)
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

	if rayTmin >= root || root >= rayTmax {
		root = (h + sqrtd) / a
		if rayTmin >= root || root >= rayTmax {
			return false
		}
	}

	hitRecord.Time = root
	hitRecord.Point = ray.At(hitRecord.Time)

	normal := mat.NewVecDense(3, nil)
	normal.SubVec(hitRecord.Point, s.center)
	// Make it unit
	normal.ScaleVec(1/s.radius, normal)
	hitRecord.SetFaceNormal(ray, normal)

	return true
}
