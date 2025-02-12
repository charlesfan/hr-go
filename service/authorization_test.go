package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/charlesfan/hr-go/config"
	"github.com/charlesfan/hr-go/service"
	"github.com/charlesfan/hr-go/test"
)

type AuthenticationServiceTestCaseSuite struct {
	env     *test.Env
	service service.AuthenticationServicer
}

func setupAuthenticationServiceTestCase(t *testing.T) (AuthenticationServiceTestCaseSuite, func(t *testing.T)) {
	s := AuthenticationServiceTestCaseSuite{
		env: test.SetupEnv(t),
	}

	s.service = service.NewAuthenticationService(service.AuthenticationServiceConfig{
		TokenExpired: time.Hour * 24 * 1, // one week
		Key:          config.APP_SECRET,
	})
	service.AuthenticationService = s.service

	return s, func(t *testing.T) {
		defer s.env.Close()
	}
}

func TestAuthenticationService_CreateToken(t *testing.T) {
	s, teardownTestCase := setupAuthenticationServiceTestCase(t)
	defer teardownTestCase(t)

	tt := []struct {
		name string

		givenJwtConfig service.JwtConfig

		wantResponseCode int

		setupSubTest test.SetupSubTest
	}{
		{
			name:             "create token with user",
			wantResponseCode: service.ErrorCodeSuccess,
			setupSubTest:     test.EmptySubTest(),
			givenJwtConfig: service.JwtConfig{
				Name: "charles-test",
			},
		},
		{
			name:             "empty jwt config",
			wantResponseCode: service.ErrorCodeSuccess,
			setupSubTest:     test.EmptySubTest(),
			givenJwtConfig:   service.JwtConfig{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			tokenStr, expiredAt, errCode := s.service.CreateToken(tc.givenJwtConfig)

			assert.Equal(t, errCode, tc.wantResponseCode)
			assert.NotZero(t, tokenStr)
			assert.NotZero(t, expiredAt)
		})
	}
}

func TestAuthenticationService_Verify(t *testing.T) {
	s, teardownTestCase := setupAuthenticationServiceTestCase(t)
	defer teardownTestCase(t)

	tt := []struct {
		name string

		givenTokenStr string

		wantResponseCode int

		setupSubTest test.SetupSubTest
	}{
		{
			name:             "verify success",
			wantResponseCode: service.ErrorCodeSuccess,
			givenTokenStr:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiY2hhcmxlcy10ZXN0IiwiZXhwIjo3OTYwMTQ4MTU3fQ.bjIjpyroDjTwg2JbCjSKwoKL2d3m_oASt1hCFZmmkgg",
			setupSubTest:     test.EmptySubTest(),
		},
		{
			name:             "expired token string verify fail",
			wantResponseCode: service.ErrorCodeTokenExpired,
			givenTokenStr:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibWJhIiwiZXhwIjoxNTM0NDA5MTUwfQ.6Ex9SqYm7w7_VAyV29WwhEOejlNv1MAXT2LB0q3YHtQ",
			setupSubTest:     test.EmptySubTest(),
		},
		{
			name:             "empty token string verify fail",
			wantResponseCode: service.ErrorCodeForbidden,
			givenTokenStr:    "",
			setupSubTest:     test.EmptySubTest(),
		},
		{
			name:             "HS512 token verify fail",
			wantResponseCode: service.ErrorCodeForbidden,
			givenTokenStr:    "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEuMzAwODE5MzhlKzA5LCJodHRwOi8vZXhhbXBsZS5jb20vaXNfcm9vdCI6dHJ1ZSwiaXNzIjoiam9lIn0.CN7YijRX6Aw1n2jyI2Id1w90ja-DEMYiWixhYCyHnrZ1VfJRaFQz1bEbjjA5Fn4CLYaUG432dEYmSbS4Saokmw",
			setupSubTest:     test.EmptySubTest(),
		},
		{
			name:             "Unexpected signing method: RS256",
			wantResponseCode: service.ErrorCodeForbidden,
			givenTokenStr:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InRlc3RAaG90bWFpbC5jb20iLCJOYW1lIjoidGVzdE5hbWUifQ.fgd1h4LB1zzAiPFLKMOJrQu12rTLeXBDKHdnqiNc04NRn-1v7cHEQpDNawvScMIGrcQLbZo6WrldZQT9ImYWpUyy3CcD2uMO95I5PN6aXOSPb26nNGQpmIi1HNZrq5359hKZ6BWEJnW9iTg7RgmMvZGmIqlGLsOY2a6UiiwBsI0",
			setupSubTest:     test.EmptySubTest(),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			_, errCode := s.service.Verify(tc.givenTokenStr)

			assert.Equal(t, errCode, tc.wantResponseCode)
		})
	}
}
