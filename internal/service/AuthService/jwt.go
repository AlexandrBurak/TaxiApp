package AuthService

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zxcghoulhunter/InnoTaxi/internal/config"
	"github.com/zxcghoulhunter/InnoTaxi/internal/model"
)

type Claims struct {
	Username string
	jwt.StandardClaims
}

func SignIn(c context.Context, login model.Login) (string, time.Time, error) {
	cfg, err := config.GetAppCfg()
	if err != nil {
		return "", time.Now(), err
	}
	duration, err := time.ParseDuration(cfg.JWT_EXPIRATION_TIME)
	if err != nil {
		return "", time.Now(), err
	}
	expirationTime := time.Now().Add(duration)

	claims := &Claims{
		Username: login.Phone,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.SECRET))
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expirationTime, nil
}
