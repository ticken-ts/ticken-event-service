package utils

import "github.com/coreos/go-oidc"

type JWTClaims struct {
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
}

func ParseJWTClaims(jwt any) (*JWTClaims, error) {
	jwtCasted := jwt.(*oidc.IDToken)

	var claims JWTClaims
	err := jwtCasted.Claims(&claims)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func GetUserIDFromJWT(jwt any) string {
	return jwt.(*oidc.IDToken).Subject
}
