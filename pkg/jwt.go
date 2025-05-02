package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
)

type Claims struct {
	Id   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// membuat konstructor jwt
func NewClaims(id int, role string) *Claims {
	return &Claims{
		Id: id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
			// jwt ini akan aktif selama 5 menit kedepan
		},
	}
}

// membuat token dengan generate
func (c *Claims) GenerateToken() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("secret not provided")
	}

	// buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// tandatangan token
	return token.SignedString([]byte(jwtSecret))
}

// melakukan verifikasi
func (c *Claims) VerifyToken(token string) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("secret not provided")
	}

	// melakukan parsing token
	parsedToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (interface{}, error) {
		// fungsi callback yang digunakan olehWithClaims untuk mengambil secret
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return err
	}

	if !parsedToken.Valid { 
		return errors.New("expired token")
	}

	return nil
}