package screening

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lcolman/fabrikam-auth-poc/internal/auth"
	"github.com/lcolman/fabrikam-auth-poc/internal/server"
)

func GetRoute(routeConfig ScreeningRouteConfig) server.Route {
	return server.Route{
		Method:  http.MethodGet,
		Path:    "/screening",
		Handler: GetHandler(routeConfig.ScreeningService(), routeConfig.Logger()),
	}
}

func GetHandler(service Service, logger *log.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		logger.Printf("INFO GET %s request handler started", req.URL.Path)
		// retrieve the authenticated subject from request context
		subject, err := auth.Subject(req.Context())
		if err != nil {
			logger.Printf("INFO rejecting as unathorized due to missing subject")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Printf("DEBUG found subject for GET request: %s", subject)
		// validate the request's OAuth scopes are authorized
		scopes, err := auth.Scopes(req.Context())
		if err != nil {
			logger.Printf("INFO rejecting as unathorized due to missing request scopes")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, hasActionAuthorization := scopes["list:screening"]
		if !hasActionAuthorization {
			logger.Printf("INFO rejecting as forbidden due to authenticated user not having scope %s", "list:screening")
			res.WriteHeader(http.StatusForbidden)
			return
		}
		// make service call to list screenings
		result, err := service.List(req.Context(), NewIdent(subject))
		if err != nil {
			logger.Printf("ERROR internal server error encountered when listing screenings: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		// format JSON output and marshal for response
		content, err := json.Marshal(result)
		if err != nil {
			logger.Printf("ERROR internal server error encountered when marshalling JSON response: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(content)
		logger.Printf("INFO GET %s request succeeded: %v %s", req.URL.Path, http.StatusOK, "OK")
	}
}
