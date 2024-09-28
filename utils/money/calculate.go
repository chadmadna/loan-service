package money

import (
	"fmt"
	"math"
	"strconv"
)

const epsilon = 1.e-9

func CalculateROI(principalStr string, annualInterestRate float64, term int) (roi string, totalInterest string) {
	principalFloat, _ := strconv.ParseFloat(principalStr, 64)

	loanInterestRate := annualInterestRate * (float64(term) / 12.0)

	totalInterestFloat := principalFloat * loanInterestRate
	roiFloat := totalInterestFloat / principalFloat * 100

	return fmt.Sprintf("%.2f", roiFloat), fmt.Sprintf("%f", toFixed(totalInterestFloat, 12))
}

// mathy stuff for comparing floats
func NearlyEqual(a, b float64) bool {
	if a == b {
		return true
	}

	diff := math.Abs(a - b)
	if a == 0.0 || b == 0.0 || diff < math.SmallestNonzeroFloat64 {
		return diff < epsilon*math.SmallestNonzeroFloat64
	}

	return diff/(math.Abs(a)+math.Abs(b)) < epsilon
}

// https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
