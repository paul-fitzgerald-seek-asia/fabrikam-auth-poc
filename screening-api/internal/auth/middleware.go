package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type middlewareKey int

const (
	subjectContextKey middlewareKey = iota
	scopesContextKey  middlewareKey = iota
)

func Subject(context context.Context) (string, error) {
	subject, present := context.Value(subjectContextKey).(string)
	if !present {
		return "", fmt.Errorf("no subject in context")
	}
	return subject, nil
}

func Scopes(context context.Context) (map[string]interface{}, error) {
	value := context.Value(scopesContextKey)
	if value == nil {
		return nil, fmt.Errorf("no scopes found in context")
	}
	scopes, correct := value.(map[string]interface{})
	if !correct {
		return nil, fmt.Errorf("type mismatch in context scopes value")
	}
	return scopes, nil
}

func Middleware(config AuthMiddlewareConfig, logger *log.Logger) mux.MiddlewareFunc {
	keySet := NewJSONWebKeySetCache(config, logger)
	tokenValidator := NewJSONWebTokenValidator(config, keySet, logger)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			logger.Printf("INFO %s request to %s from %s\n", req.Method, req.URL.Path, req.RemoteAddr)
			if req.Method == http.MethodOptions {
				next.ServeHTTP(writer, req)
				return
			}
			authHeader := req.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				next.ServeHTTP(writer, req)
			}
			rawToken := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := tokenValidator.ValidateTokenString(rawToken)
			if err != nil {
				logger.Printf("ERROR Error during token validation: %v", err)
				next.ServeHTTP(writer, req)
				return
			}
			if token == nil {
				next.ServeHTTP(writer, req)
				return
			}
			logger.Printf("DEBUG Valid bearer token received from client request")
			tokenSubject, subjectPresent := token["sub"]
			tokenScopes, scopesPresent := token["permissions"]
			if !subjectPresent || !scopesPresent {
				logger.Printf("DEBUG Rejecting token for context annotation: missing required claims")
				next.ServeHTTP(writer, req)
				return
			}
			subject, subCastOk := tokenSubject.(string)
			rawScopes, scopeCastOk := tokenScopes.([]interface{})
			if !subCastOk || !scopeCastOk {
				logger.Printf("DEBUG Aborting request context annotation: unable to type cast claims")
				next.ServeHTTP(writer, req)
				return
			}
			// tweak the scopes data structure as Go has no native Set abstraction
			scopes := map[string]interface{}{}
			for _, scope := range rawScopes {
				scopes[scope.(string)] = nil
			}
			oldContext := req.Context()
			newContext := context.WithValue(oldContext, subjectContextKey, subject)
			newContext = context.WithValue(newContext, scopesContextKey, scopes)
			logger.Printf("DEBUG Successfully annotated request context with token audience and scopes")
			next.ServeHTTP(writer, req.WithContext(newContext))
		})
	}
}
