package auth

type AuthMiddlewareConfig interface {
	TokenAudience() string
	TokenIssuer() string
	KeySetHost() string
}
