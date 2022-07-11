package util

func RoundFloat(num float64) float64 {
	return float64(int((num+0.000001)*10000)) / 10000
}

func Min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
