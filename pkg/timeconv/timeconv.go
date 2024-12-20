package timeconv

import "time"

// ToTime converts a Unix timestamp (int64) to time.Time.
func ToTime(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// ToTimePtr converts a pointer to a Unix timestamp (int64) to a pointer to time.Time.
// Returns nil if the input is nil.
func ToTimePtr(unix *int64) *time.Time {
	if unix == nil {
		return nil
	}
	t := time.Unix(*unix, 0)
	return &t
}

// ToUnix converts time.Time to a Unix timestamp (int64).
func ToUnix(t time.Time) int64 {
	return t.Unix()
}

// ToUnixPtr converts a pointer to time.Time to a pointer to a Unix timestamp (int64).
// Returns nil if the input is nil.
func ToUnixPtr(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	unix := t.Unix()
	return &unix
}
