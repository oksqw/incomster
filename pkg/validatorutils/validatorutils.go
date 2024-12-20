package validatorutils

import "fmt"

// ValidateString returns an error if value is empty string
func ValidateString(value, placeholder string) error {
	if value == "" {
		return fmt.Errorf("%s is required", placeholder)
	}

	return nil
}

// ValidateInt returns an error if value is less than or equal to min
func ValidateInt(value, min int, placeholder string) error {
	if value <= min {
		return fmt.Errorf("invalid %s", placeholder)
	}
	return nil
}

// ValidateFloat64 returns an error if value is less than or equal to min
func ValidateFloat64(value, min float64, placeholder string) error {
	if value <= min {
		return fmt.Errorf("invalid %s", placeholder)
	}
	return nil
}
