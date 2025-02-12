package service

import (
	"time"

	"github.com/charlesfan/hr-go/config"
)

var (
	// === Repository ===
	// === Service ===
	AuthenticationService AuthenticationServicer
)

func Init() {
	// === Service ===
	// authentication service
	AuthenticationService = NewAuthenticationService(AuthenticationServiceConfig{
		TokenExpired: time.Hour * 24 * 7, // one week
		Key:          config.APP_SECRET,
	})
}
