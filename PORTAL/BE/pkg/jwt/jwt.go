package jwt

import (
	"portal_api/config"
	"portal_api/errors"
	"portal_api/helper"
	time_helper "portal_api/helper/time"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	claims                 *jwt.MapClaims
	invalidTokenErrMessage = "Invalid token"
	expiredTokenErrMessage = "Token expired"
)

func CreateClaims(data map[string]any) {
	jwtClaims := jwt.MapClaims{}
	for key, value := range data {
		jwtClaims[key] = value
	}

	// jwtClaims["exp"] = time_helper.GetTimeNow().Add(time.Hour * 2).Unix()
	expireDuration := config.Cfg.TokenExpire
	jwtClaims["exp"] = time_helper.GetTimeNow().Add(time.Second * time.Duration(expireDuration)).Unix()
	claims = &jwtClaims
}

func GenerateToken() (string, error) {
	secretKey, err := helper.GetSecretKey()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	secretKey, err := helper.GetSecretKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		var errMessage string
		if strings.Contains(err.Error(), "expired") {
			errMessage = expiredTokenErrMessage
		} else {
			errMessage = invalidTokenErrMessage
		}
		return nil, &errors.UnauthenticatedError{Message: errMessage}
	}
	if !token.Valid {
		return nil, &errors.UnauthenticatedError{Message: invalidTokenErrMessage}
	}

	jwtClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &errors.UnauthenticatedError{Message: invalidTokenErrMessage}
	}

	expUnix, ok := jwtClaims["exp"].(float64)
	if !ok {
		return nil, &errors.UnauthenticatedError{Message: invalidTokenErrMessage}
	}
	expTime := time.Unix(int64(expUnix), 0)
	if time_helper.GetTimeNow().After(expTime) {
		return nil, &errors.UnauthenticatedError{Message: expiredTokenErrMessage}
	}

	return &jwtClaims, nil
}
