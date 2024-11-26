package ray

import "raytracer/util"

type HittableList struct {
	objects []Hittable
}

func NewHittableList() HittableList {
	return HittableList{
		objects: make([]Hittable, 0),
	}
}

func (h *HittableList) Add(obj Hittable) {
	h.objects = append(h.objects, obj)
}

func (h *HittableList) Clear() {
	h.objects = nil
}

func (h *HittableList) Hit(ray *Ray, interval util.Interval, hitRecord *HitRecord) bool {
	tempRec := NewHitRecord()
	hitAnything := false
	closestSoFar := interval.Max

	for _, object := range h.objects {
		hit := object.Hit(ray, util.NewInterval(interval.Min, closestSoFar), &tempRec)

		if hit {
			hitAnything = true
			closestSoFar = tempRec.Time
			*hitRecord = tempRec
		}
	}

	return hitAnything
}
