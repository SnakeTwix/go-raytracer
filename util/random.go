package util

import "math/rand/v2"

func RandomF64() float64 {
	return rand.Float64()
}

func RandomF64Range(min float64, max float64) float64 {
	return min + (max-min)*RandomF64()
}
