package material

import (
	"gonum.org/v1/gonum/mat"
	"raytracer/ray"
)

type Material interface {
	Scatter(rayIn *ray.Ray, hitRecord *ray.HitRecord, attenuation *mat.VecDense, scatted *ray.Ray) bool
}
