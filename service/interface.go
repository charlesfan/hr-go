package service

type AuthenticationServicer interface {
	CreateToken(cfg JwtConfig) (tokenStr string, expiresAt int64, errCode int)
	Verify(tokenStr string) (claims *CustomClaims, errCode int)
}

type EmployeeServicer interface {
	LoginByEmailPassword(string, string) (string, int64, int)
}
