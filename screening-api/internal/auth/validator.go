package auth

import (
	"errors"
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

type JSONWebTokenValidator interface {
	ValidateTokenString(string) (map[string]interface{}, error)
}

type jwtValidator struct {
	expectedAlgorithm *jwt.SigningMethodRSA
	expectedAudience  string
	expectedIssuer    string
	keySet            JSONWebKeySetCache
	logger            *log.Logger
}

func NewJSONWebTokenValidator(config AuthMiddlewareConfig, keySet JSONWebKeySetCache, logger *log.Logger) JSONWebTokenValidator {
	return &jwtValidator{
		expectedAlgorithm: jwt.SigningMethodRS256, //hard-coded for now
		expectedAudience:  config.TokenAudience(),
		expectedIssuer:    config.TokenIssuer(),
		keySet:            keySet,
		logger:            logger,
	}
}

func ExtractTokenHeader(token *jwt.Token, headerKey string) (string, error) {
	if token.Header == nil {
		return "", errors.New("Invalid token structure: missing header")
	}
	header, notMissing := token.Header[headerKey]
	if isMissing := !notMissing; isMissing {
		return "", fmt.Errorf("Token missing required header: %s", headerKey)
	}
	headerString, wasString := header.(string)
	if !wasString {
		return "", errors.New("Invalid data type found in token header")
	}
	return headerString, nil
}

func (v *jwtValidator) ValidateAlgorithm(token *jwt.Token) error {
	tokenAlgorithm, err := ExtractTokenHeader(token, "alg")
	if err != nil {
		return err
	}
	if v.expectedAlgorithm.Alg() != tokenAlgorithm {
		return errors.New("Illegal token signature algorithm specified")
	}
	return nil
}

// handleMultipleTokenAudiences modifies the given claims map in-place to ensure aud is always an array
func handleMultipleTokenAudiences(claims *jwt.MapClaims) {
	var newAudArray []string
	singleAudience, found := (*claims)["aud"].(string)
	if !found {
		multiAudience, _ := (*claims)["aud"].([]interface{})
		for _, audClaim := range multiAudience {
			newAudArray = append(newAudArray, audClaim.(string))
		}
	} else {
		newAudArray = []string{singleAudience}
	}
	(*claims)["aud"] = newAudArray
}

func (v *jwtValidator) ValidateClaims(token *jwt.Token) (interface{}, error) {
	claimsMap, correct := token.Claims.(jwt.MapClaims)
	if !correct {
		return nil, errors.New("Invalid token claims structure")
	}
	expirationErr := claimsMap.Valid()
	if expirationErr != nil {
		return nil, errors.New("Invalid token expiry")
	}
	//validAud := claimsMap.VerifyAudience(v.expectedAudience, true)  // chokes when aud is an array
	// apparantly auth0 issues tokens with multiple aud claims
	handleMultipleTokenAudiences(&claimsMap) // so I needed this hack
	validAud := false
	for _, tokenAud := range claimsMap["aud"].([]string) {
		if tokenAud == v.expectedAudience {
			validAud = true
			break
		}
	}
	// </hack>
	if !validAud {
		return nil, errors.New("Invalid token audience")
	}
	validIss := claimsMap.VerifyIssuer(v.expectedIssuer, true)
	if !validIss {
		return token, errors.New("Invalid token issuer")
	}
	keyID, err := ExtractTokenHeader(token, "kid")
	if err != nil {
		return nil, err
	}
	signingKey, err := v.keySet.FindPEMKey(keyID)
	if err != nil {
		return nil, err
	}
	return signingKey, nil
}

func (v *jwtValidator) ValidateToken(token *jwt.Token) (interface{}, error) {
	err := v.ValidateAlgorithm(token)
	if err != nil {
		return nil, err
	}
	return v.ValidateClaims(token)
}

func (v *jwtValidator) ValidateTokenString(raw string) (map[string]interface{}, error) {
	token, err := jwt.Parse(raw, v.ValidateToken)
	if err != nil {
		return nil, err
	}
	return (token.Claims.(jwt.MapClaims)), nil
}
