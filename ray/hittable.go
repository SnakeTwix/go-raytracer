package ray

import (
	"gonum.org/v1/gonum/mat"
	"raytracer/util"
)

type HitRecord struct {
	Normal    *mat.VecDense
	Point     *mat.VecDense
	Time      float64
	FrontFace bool
}

func NewHitRecord() HitRecord {
	return HitRecord{}
}

func (h *HitRecord) SetFaceNormal(ray *Ray, normal *mat.VecDense) {
	h.Normal = normal

	if mat.Dot(ray.Direction, normal) < 0 {
		h.FrontFace = true
	} else {
		h.FrontFace = false
		h.Normal.ScaleVec(-1, h.Normal)
	}
}

type Hittable interface {
	Hit(ray *Ray, interval util.Interval, hitRecord *HitRecord) bool
}
