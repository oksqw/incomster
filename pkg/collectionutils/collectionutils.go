package collectionutils

// ToInterface converts a slice of any type to a slice of interfaces.
// It takes a generic slice as input and returns a slice where each element
// is wrapped as an empty interface.
//
// Parameters:
// - input []T: A slice of any type.
//
// Returns:
// - []interface{}: A slice of interface{} containing all elements of the input slice.
//
// Example:
// Given an input slice []int{1, 2, 3}, the result will be []interface{}{1, 2, 3}.
func ToInterface[T any](input []T) []interface{} {
	result := make([]interface{}, len(input))
	for i, v := range input {
		result[i] = v
	}
	return result
}
