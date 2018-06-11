package math

import (
	m "math"
)

func Sine(n int, amp, base, period float64) []float64 {
	amp, base = amp - 0.001, base - 0.001
	res := make([]float64, n)
	for i := 0; i < n; i++ {
		res[i] = base + amp * m.Cos(period * m.Pi * float64(i) / float64(n))
	}
	return res
}
