package helper

import (
	"fmt"
	"golang-todo/internal/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func GetDateNow() time.Time {
	// Load the IANA time zone for WIB.
	wibLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Time{} // Return zero time and the error
	}

	// Get the current time and convert it to the WIB location.
	return time.Now().In(wibLocation)
}

func GenerateJWTToken(data model.UserCustomClaims, log *logrus.Logger) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)
	secretKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		log.Errorf("got error when sign a jwt token %v", err)
	}

	return tokenString
}

func ParseJWTToken(log *logrus.Logger, token string, claims jwt.Claims, options ...jwt.ParserOption) (*model.UserCustomClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	cb := func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secretKey), nil
	}
	parsedToken, err := jwt.ParseWithClaims(token, claims, cb, options...)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	if claims, ok := parsedToken.Claims.(*model.UserCustomClaims); ok && parsedToken.Valid {
		log.Info("Token valid")
		return claims, nil
	} else {
		log.Info("Token invalid")
		return nil, err
	}
}
