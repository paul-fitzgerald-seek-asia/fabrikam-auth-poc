package auth

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type JSONWebKeySetCache interface {
	FindPEMKey(keyID string) (*rsa.PublicKey, error)
	FindCertString(keyID string) (string, error)
}

type keySetCache struct {
	expectedAudience string
	expectedIssuer   string
	hostDomain       string
	logger           *log.Logger
	keyMap           struct {
		Keys []struct {
			Kty string   `json:"kty"`
			Kid string   `json:"kid"`
			Use string   `json:"use"`
			N   string   `json:"n"`
			E   string   `json:"e"`
			X5c []string `json:"x5c"`
		} `json:"keys"`
	}
}

func NewJSONWebKeySetCache(config AuthMiddlewareConfig, logger *log.Logger) JSONWebKeySetCache {
	cache := keySetCache{
		hostDomain: config.KeySetHost(),
		logger:     logger,
	}
	err := cache.prefetchCertificates()
	if err != nil {
		panic(err)
	}
	return &cache
}

func (ks *keySetCache) prefetchCertificates() error {
	ks.logger.Printf("INFO Fetching JWKS for domain %s", ks.hostDomain)
	resp, err := http.Get(fmt.Sprintf("https://%s/.well-known/jwks.json", ks.hostDomain))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ks.keyMap)
	if err != nil {
		ks.logger.Printf("ERROR Unable to parse JWKS response for domain %s: %v", ks.hostDomain, err)
		return err
	}
	ks.logger.Printf("INFO Successfully loaded %v JWKS certificates for domain %s", len(ks.keyMap.Keys), ks.hostDomain)
	return nil
}

func (ks *keySetCache) FindPEMKey(keyID string) (*rsa.PublicKey, error) {
	certString, err := ks.FindCertString(keyID)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(certString))
	if err != nil {
		ks.logger.Printf("ERROR Incorrect format encountered while attempting to parse PEM key")
		return nil, err
	}
	return publicKey, nil
}

func (ks *keySetCache) FindCertString(keyID string) (string, error) {
	for _, key := range ks.keyMap.Keys {
		if keyID == key.Kid {
			ks.logger.Printf("DEBUG Successfully found matching certificate key id for kid=\"%s\"", key.Kid)
			return fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", key.X5c[0]), nil
		}
	}
	ks.logger.Printf("ERROR Unable to find mathing key in JWK set for domain %s", ks.hostDomain)
	return "", errors.New("unable to find matching key in configured set")
}
