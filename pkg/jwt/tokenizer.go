package jwt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cristalhq/jwt/v4"
)

type Tokenizer struct {
	signer   jwt.Signer
	verifier jwt.Verifier
	duration time.Duration
}

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func New(secret string, duration time.Duration) (*Tokenizer, error) {
	key := []byte(secret)

	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer: %w", err)
	}

	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		return nil, fmt.Errorf("failed to create verifier: %w", err)
	}

	return &Tokenizer{
		signer:   signer,
		verifier: verifier,
		duration: duration,
	}, nil
}

// Generate створює токен
func (t *Tokenizer) Generate(userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "incomster",
		},
	}

	token, err := jwt.NewBuilder(t.signer).Build(&claims)
	if err != nil {
		return "", fmt.Errorf("failed to build token: %w", err)
	}

	return token.String(), nil
}

// Parse перевіряє токен і повертає claims
func (t *Tokenizer) Parse(encoded string) (*Claims, error) {
	token, err := jwt.Parse([]byte(encoded), t.verifier)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	var claims Claims
	err = json.Unmarshal(token.Claims(), &claims)
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	//if !claims.IsValidAt(time.Now()) {
	//	return nil, fmt.Errorf("token is not valid at the current time")
	//}

	return &claims, nil
}

//package jwt
//
//import (
//	"github.com/golang-jwt/jwt/v5"
//	"time"
//)
//
//type Tokenizer struct {
//	secret   []byte
//	duration time.Duration
//}
//
//func New(secret string, duration time.Duration) *Tokenizer {
//	return &Tokenizer{
//		secret:   []byte(secret),
//		duration: duration,
//	}
//}
//
//type Claims struct {
//	UserID int    `json:"user_id"`
//	Role   string `json:"role"`
//	jwt.RegisteredClaims
//}
//
//// Generate creates a JWT token for the user
//func (t *Tokenizer) Generate(userID int, role string) (string, error) {
//	claims := &Claims{
//		UserID: userID,
//		Role:   role,
//		RegisteredClaims: jwt.RegisteredClaims{
//			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.duration)),
//			IssuedAt:  jwt.NewNumericDate(time.Now()),
//			Issuer:    "incomster",
//		},
//	}
//
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	return token.SignedString(t.secret)
//}
//
//// Parse checks the token and returns Claims
//func (t *Tokenizer) Parse(hash string) (*Claims, error) {
//	token, err := jwt.ParseWithClaims(hash, &Claims{}, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, jwt.ErrInvalidKey
//		}
//		return t.secret, nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	claims, ok := token.Claims.(*Claims)
//	if !ok || !token.Valid {
//		return nil, jwt.ErrTokenMalformed
//	}
//
//	return claims, nil
//}
