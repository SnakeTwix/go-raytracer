package hit

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

func (h *HittableList) Hit(ray *Ray, rayTmin float64, rayTmax float64, hitRecord *Record) bool {
	tempRec := NewRecord()
	hitAnything := false
	closestSoFar := rayTmax

	for _, object := range h.objects {
		hit := object.Hit(ray, rayTmin, closestSoFar, &tempRec)

		if hit {
			hitAnything = true
			closestSoFar = tempRec.Time
			*hitRecord = tempRec
		}
	}

	return hitAnything
}
