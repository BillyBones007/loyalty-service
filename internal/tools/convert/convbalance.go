package convert

// Convertation balance integer to float64
func ConvToFloatBalance(intBalance int64) float64 {
	floatBalance := float64(intBalance) / 100.0
	return floatBalance
}

// Convertation balance float64 to integer
func ConvToIntBalance(floatBalance float64) int64 {
	intBalance := floatBalance * 100
	return int64(intBalance)
}
