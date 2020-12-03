package config

import (
	"fmt"
	"os"
	"strconv"
)

type ApplicationConfig interface {
	ServerPort() int
	BasePath() string
	LoggingName() string
	TokenAudience() string
	TokenIssuer() string
	KeySetHost() string
}

type appConfig struct {
	port       int
	apiVersion int
	apiName    string
	audience   string
	issuer     string
	JWKset     string
}

func (c *appConfig) ServerPort() int {
	return c.port
}

func (c *appConfig) BasePath() string {
	return fmt.Sprintf("/v%v/%s", c.apiVersion, c.apiName)
}

func (c *appConfig) LoggingName() string {
	return c.apiName
}

func (c *appConfig) TokenAudience() string {
	return c.audience
}

func (c *appConfig) TokenIssuer() string {
	return c.issuer
}

func (c *appConfig) KeySetHost() string {
	return c.JWKset
}

func fromEnvOrDefaultInt(varName string, defaultValue int) int {
	num := defaultValue
	strValue := os.Getenv(varName)
	if strValue != "" {
		var err error
		num, err = strconv.Atoi(strValue)
		if err != nil {
			return defaultValue
		}
	}
	return num
}

func fromEnvOrDefaultStr(varName string, defaultValue string) string {
	str := defaultValue
	envValue := os.Getenv(varName)
	if envValue != "" {
		str = envValue
	}
	return str
}

func NewApplicationConfig() ApplicationConfig {
	return &appConfig{
		port:       fromEnvOrDefaultInt("SCREENING_API_PORT", 8080),
		apiVersion: fromEnvOrDefaultInt("SCREENING_API_VERSION", 1),
		apiName:    fromEnvOrDefaultStr("SCREENING_API_NAME", "screening"),
		audience:   fromEnvOrDefaultStr("SCREENING_API_JWT_AUDIENCE", "http://localhost:8080/v1/screening/"),
		issuer:     fromEnvOrDefaultStr("SCREENING_API_JWT_ISSUER", "https://lawrence-seek-idp-dev.au.auth0.com/"),
		JWKset:     fromEnvOrDefaultStr("SCREENING_API_JWKS_HOST", "lawrence-seek-idp-dev.au.auth0.com"),
	}
}
