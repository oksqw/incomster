package ternary

// Func returns t if c is true, otherwise returns f.
func Func[T any](c bool, t, f T) T {
	if c {
		return t
	}
	return f
}
