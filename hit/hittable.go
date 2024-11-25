package hit

import (
	"gonum.org/v1/gonum/mat"
)

type Record struct {
	Normal    *mat.VecDense
	Point     *mat.VecDense
	Time      float64
	FrontFace bool
}

func NewRecord() Record {
	return Record{}
}

func (h *Record) SetFaceNormal(ray *Ray, normal *mat.VecDense) {
	h.Normal = normal

	if mat.Dot(ray.Direction, normal) < 0 {
		h.FrontFace = true
	} else {
		h.FrontFace = false
		h.Normal.ScaleVec(-1, h.Normal)
	}
}

type Hittable interface {
	Hit(ray *Ray, rayTmin float64, rayTmax float64, hitRecord *Record) bool
}
