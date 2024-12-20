package mapping

func Map[T any, R any](input []T, convert func(T) R) []R {
	output := make([]R, len(input))
	for i, v := range input {
		output[i] = convert(v)
	}
	return output
}
