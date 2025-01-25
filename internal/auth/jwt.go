package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuthenticator struct {
	secret   string
	audience string
	iss      string
}

func NewJWTAuthenticator(secret, audience, iss string) JWTAuthenticator {
	return JWTAuthenticator{
		secret:   secret,
		audience: audience,
		iss:      iss,
	}
}

func (a *JWTAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", nil
	}
}

func (a *JWTAuthenticator) VerifyToken(token string) *jwt.Token {

}
