package service

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

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
		return token.Claims.(*CustomClaims), ErrorCodeSuccess
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

func NewAuthenticationService(cfg AuthenticationServiceConfig) AuthenticationServicer {
	return &authenticationService{
		tokenExpired: cfg.TokenExpired,
		key:          cfg.Key,
	}
}
