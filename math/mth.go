// Package math provides additional mathematical functions extending the standard math library.
package math

import (
	"math"
	"sort"
)

// Constants not available in the standard math library
const (
	GoldenRatio = 1.6180339887498948482
	E10         = 22026.465794806716517
	Ln10        = 2.3025850929940456840
	Sqrt3       = 1.7320508075688772935
	Sqrt5       = 2.2360679774997896964
)

// Mean calculates the arithmetic mean (average) of a slice of float64 values.
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

// Median calculates the median of a slice of float64 values.
func Median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy to avoid modifying the original slice
	valuesCopy := make([]float64, len(values))
	copy(valuesCopy, values)
	sort.Float64s(valuesCopy)

	// Get the median value
	n := len(valuesCopy)
	if n%2 == 0 {
		return (valuesCopy[n/2-1] + valuesCopy[n/2]) / 2
	}
	return valuesCopy[n/2]
}

// StandardDeviation calculates the population standard deviation.
func StandardDeviation(values []float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	m := Mean(values)
	var sum float64
	for _, v := range values {
		sum += math.Pow(v-m, 2)
	}
	return math.Sqrt(sum / float64(len(values)))
}

// Factorial calculates the factorial of n.
func Factorial(n int) int {
	if n <= 0 {
		return 1
	}
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// Combination calculates the number of ways to choose k items from n items without repetition and without order.
func Combination(n, k int) int {
	if k > n || k < 0 {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	result := 1
	for i := 1; i <= k; i++ {
		result *= (n - (k - i))
		result /= i
	}
	return result
}

// CompoundInterest calculates the compound interest.
// p = principal, r = annual interest rate, n = number of times compounded per year, t = time in years
func CompoundInterest(p, r, n, t float64) float64 {
	return p * math.Pow(1+(r/n), n*t)
}

// DegreesToRadians converts degrees to radians.
func DegreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// RadiansToDegrees converts radians to degrees.
func RadiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// IsPrime determines if n is a prime number.
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	// Check all numbers of form 6k ± 1 up to sqrt(n)
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// GCD finds the greatest common divisor of two integers.
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM finds the least common multiple of two integers.
func LCM(a, b int) int {
	return a / GCD(a, b) * b
}
