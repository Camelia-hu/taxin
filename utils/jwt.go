package utils

import (
	"errors"
	"github.com/Camelia-hu/taxin/model"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"time"
)

func DeliverTokenByRPCService(userId string) (string, string, error) {
	myAccessClaims := model.MyClaims{
		Id: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * 60 * time.Minute)),
		},
	}
	myRefreshClaims := model.MyClaims{
		Id: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	}
	AccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, myAccessClaims).SignedString([]byte(os.Getenv("JWT_KEY")))
	RefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, myRefreshClaims).SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", "", err
	}
	return AccessToken, RefreshToken, nil
}

func VerifyTokenByRPCService(Token string) (err error) {
	myClaims := new(model.MyClaims)
	token, err := jwt.ParseWithClaims(Token, myClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		//span.SetStatus(codes.Error, "token parse err")
		log.Println("1 ", err)
		return err
	}
	if !token.Valid {
		//span.SetStatus(codes.Error, "token expired")
		return errors.New("token 过期喵～")
	}
	return nil
}

func ReFreshTokenByRPCService(RefreshToken string) (NewAccessToken string, NewRefreshToken string, err error) {
	//span := trace.SpanFromContext(s.ctx)
	var myClaims model.MyClaims
	refreshToken, err := jwt.ParseWithClaims(RefreshToken, &myClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		//span.SetStatus(codes.Error, err.Error())
		return "", "", err
	}
	if !refreshToken.Valid {
		//span.SetStatus(codes.Error, "refreshToken 过期 too 喵～")
		return "", "", errors.New("refreshToken 过期 too 喵～")
	}
	newrefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, model.MyClaims{
		Id: myClaims.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		},
	}).SignedString([]byte(os.Getenv("JWT_KEY")))
	newaccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, model.MyClaims{
		Id: myClaims.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		},
	}).SignedString([]byte(os.Getenv("JWT_KEY")))

	return newaccessToken, newrefreshToken, nil
}
