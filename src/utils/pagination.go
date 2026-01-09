package utils

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// NormalizeLimit ensures limit is within valid bounds
func NormalizeLimit(limit int) int {
	if limit <= 0 || limit > MaxLimit {
		return DefaultLimit
	}
	return limit
}

// NormalizeOffset ensures offset is non-negative
func NormalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
