package passwordutil

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// Hash hashes the password using bcrypt
func Hash(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// Compare compares the entered password with the hash
func Compare(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}

func Validate(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// TODO: Enable in release mode?
	// Отключено из-за ненадобности на данный момент.
	// И не будет включено, т.к. это тестовый проект и мне лень печатать постоянно такие пароли.
	//
	//var hasUpper, hasLower, hasDigit, hasSpecial bool
	//
	//for _, char := range password {
	//	switch {
	//	case unicode.IsUpper(char):
	//		hasUpper = true
	//	case unicode.IsLower(char):
	//		hasLower = true
	//	case unicode.IsDigit(char):
	//		hasDigit = true
	//	case unicode.IsPunct(char) || unicode.IsSymbol(char):
	//		hasSpecial = true
	//	}
	//
	//	if hasUpper && hasLower && hasDigit && hasSpecial {
	//		break
	//	}
	//}
	//
	//if !hasUpper {
	//	return errors.New("password must contain at least one uppercase letter")
	//}
	//if !hasLower {
	//	return errors.New("password must contain at least one lowercase letter")
	//}
	//if !hasDigit {
	//	return errors.New("password must contain at least one digit")
	//}
	//if !hasSpecial {
	//	return errors.New("password must contain at least one special character")
	//}

	return nil
}
