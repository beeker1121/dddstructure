package utils

import "math"

// round defines the rounding algorithms.
type round string

const (
	Round   round = "round"
	Floor   round = "floor"
	Ceil    round = "ceil"
	Bankers round = "bankers"
)

// PercentageFromInt will give you a percentage to the exact precision that you
// want based on the given fraction. For example, a fraction of 0 will return a
// number with no decimal points, rounded to the nearest whole number.
func PercentageFromInt(amt int, percentage float64, fraction int, round round) float64 {
	// Calculate percentage.
	val := float64(amt) * percentage
	val = val / 100

	// Remove potential rounding errors by moving decimal
	// two places past desired fraction and truncating.
	val = math.Trunc(val*math.Pow10(fraction+2)) / math.Pow10(fraction+2)

	// Handle rounding.
	switch round {
	case Round:
		val = math.Round(val*math.Pow10(fraction)) / math.Pow10(fraction)
	case Floor:
		val = math.Floor(val*math.Pow10(fraction)) / math.Pow10(fraction)
	case Ceil:
		val = math.Ceil(val*math.Pow10(fraction)) / math.Pow10(fraction)
	case Bankers:
		val = math.RoundToEven(val*math.Pow10(fraction)) / math.Pow10(fraction)
	default:
		val = math.Round(val*math.Pow10(fraction)) / math.Pow10(fraction)
	}

	return val
}
