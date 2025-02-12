package service

type AuthenticationServicer interface {
	CreateToken(cfg JwtConfig) (tokenStr string, expiresAt int64, errCode int)
	Verify(tokenStr string) (claims *CustomClaims, errCode int)
}
