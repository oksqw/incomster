package usernameutil

import "errors"

func Validate(input string) error {
	if len(input) < 3 || len(input) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	if input[0] == '_' || input[len(input)-1] == '_' {
		return errors.New("username cannot start or end with an underscore")
	}

	for _, char := range input {
		if !(char >= 'a' && char <= 'z') &&
			!(char >= 'A' && char <= 'Z') &&
			!(char >= '0' && char <= '9') &&
			char != '_' {
			return errors.New("username can only contain letters, numbers, and underscores")
		}
	}

	return nil
}
