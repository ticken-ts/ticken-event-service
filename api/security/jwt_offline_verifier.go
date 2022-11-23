package security

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt"
)

type JWTOfflineVerifier struct {
	key *rsa.PrivateKey
}

func NewJWTOfflineVerifier(key *rsa.PrivateKey) *JWTOfflineVerifier {
	return &JWTOfflineVerifier{key: key}
}

func (jwtVerifier *JWTOfflineVerifier) Verify(rawJWT string) (*JWT, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// here we are assuming that the JWT
		// is generated with the correct key
		return &jwtVerifier.key.PublicKey, nil
	}

	claims := new(Claims)
	_, err := jwt.ParseWithClaims(rawJWT, claims, keyFunc)
	if err != nil {
		return nil, err
	}

	token := &JWT{
		Subject:  claims.Subject,
		Email:    claims.Email,
		Username: claims.PreferredUsername,
	}

	return token, nil
}
