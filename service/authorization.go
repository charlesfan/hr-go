package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/charlesfan/hr-go/repository/cache"
	"github.com/charlesfan/hr-go/utils/log"
)

type JwtConfig struct {
	Name string `json:"name"`
}

type CustomClaims struct {
	JwtConfig
	jwt.StandardClaims
}

type authenticationService struct {
	tokenExpired time.Duration
	key          string
	cache        cache.ICache
}

type AuthenticationServiceConfig struct {
	TokenExpired time.Duration
	Key          string
}

// authentication service CreateToken method
func (s *authenticationService) CreateToken(cfg JwtConfig) (string, int64, int) {
	expiresAt := time.Now().Add(s.tokenExpired).Unix()
	claims := &CustomClaims{
		JwtConfig{
			Name: cfg.Name,
		},
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.key))
	if err != nil {
		log.Error("AuthenticationService.CreateToken fail, err => ", err)
		return "", expiresAt, ErrorCodeTokenCreateFail
	}
	return tokenStr, expiresAt, ErrorCodeSuccess
}

func (s *authenticationService) Verify(tokenStr string) (*CustomClaims, int) {
	var cc *CustomClaims

	ok, err := s.cache.BindJSON(context.Background(), tokenStr, cc)
	if ok && err == nil {
		return cc, ErrorCodeSuccess
	}
	// parser token string
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.key), nil
	})

	if token == nil && err != nil {
		log.Error("AuthenticationService.jwt.ParseWithClaims parse fail => ", err)
		return nil, ErrorCodeForbidden
	}

	if token.Valid {
		re := token.Claims.(*CustomClaims)
		go func(key string, val *CustomClaims) {
			b, err := json.Marshal(val)
			if err != nil {
				log.Error(err)
				return
			}
			if err := s.cache.Set(context.Background(), key, string(b), time.Hour*2); err != nil {
				log.Error(err)
			}
		}(tokenStr, re)
		return re, ErrorCodeSuccess
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			log.Error("AuthenticationService.Verify parse token fail => That's not even a token")
			return nil, ErrorCodeForbidden
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			log.Error("AuthenticationService.Verify parse token fail => Token is either expired or not active yet ", err)
			return nil, ErrorCodeTokenExpired
		} else {
			log.Error("AuthenticationService.Verify parse token fail => Couldn't handle this token: ", err)
			return nil, ErrorCodeForbidden
		}
	} else {
		log.Error("AuthenticationService.Verify parse token fail => Couldn't handle this token: ", err)
		return nil, ErrorCodeForbidden
	}
}

func NewAuthenticationService(cfg AuthenticationServiceConfig, c cache.ICache) AuthenticationServicer {
	return &authenticationService{
		tokenExpired: cfg.TokenExpired,
		key:          cfg.Key,
		cache:        c,
	}
}
