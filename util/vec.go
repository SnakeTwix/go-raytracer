package util

import (
	"gonum.org/v1/gonum/mat"
)

func NewZeroVector() *mat.VecDense {
	return mat.NewVecDense(3, nil)
}

func NewRandomVector() *mat.VecDense {
	return mat.NewVecDense(3, []float64{RandomF64(), RandomF64(), RandomF64()})
}

func NewRandomVectorRange(min, max float64) *mat.VecDense {
	return mat.NewVecDense(3, []float64{RandomF64Range(min, max), RandomF64Range(min, max), RandomF64Range(min, max)})
}

func NewRandomUnitVector() *mat.VecDense {
	for {
		vec := NewRandomVectorRange(-1, 1)
		lengthSquared := mat.Dot(vec, vec)

		// Floating point issues
		if 1e-160 < lengthSquared && lengthSquared <= 1 {
			MakeUnitVector(vec)
			return vec
		}
	}
}

func NewRandomUnitVectorOnHemisphere(n *mat.VecDense) *mat.VecDense {
	unitVector := NewRandomUnitVector()
	if mat.Dot(n, unitVector) < 0 {
		unitVector.ScaleVec(-1, unitVector)
	}
	return unitVector
}

// MakeUnitVector modifies the vector in place to be a unit vector
func MakeUnitVector(vec *mat.VecDense) {
	vec.ScaleVec(1./vec.Norm(2), vec)
}
