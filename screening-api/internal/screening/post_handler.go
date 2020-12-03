package screening

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lcolman/fabrikam-auth-poc/internal/auth"
	"github.com/lcolman/fabrikam-auth-poc/internal/server"
	"github.com/lcolman/fabrikam-auth-poc/pkg/definitions"
)

func PostRoute(routeConfig ScreeningRouteConfig) server.Route {
	return server.Route{
		Method:  http.MethodPost,
		Path:    "/screening",
		Handler: PostHandler(routeConfig.ScreeningService(), routeConfig.Logger()),
	}
}

func PostHandler(service Service, logger *log.Logger) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Access-Control-Allow-Origin", "*")
		logger.Printf("INFO POST %s request handler started", req.URL.Path)
		// retrieve the authenticated subject from request context
		subject, err := auth.Subject(req.Context())
		if err != nil {
			logger.Printf("INFO rejecting as unathorized due to missing subject")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.Printf("DEBUG found subject for POST request: %s", subject)
		// validate the request's OAuth scopes are authorized
		scopes, err := auth.Scopes(req.Context())
		if err != nil {
			logger.Printf("INFO rejecting as unathorized due to missing request scopes")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, hasActionAuthorization := scopes["create:screening"]
		if !hasActionAuthorization {
			logger.Printf("INFO rejecting as forbidden due to authenticated user not having scope %s", "create:screening")
			res.WriteHeader(http.StatusForbidden)
			return
		}
		// parse the submitted screening JSON content
		var newScreening definitions.Screening
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Printf("ERROR rejecting as bad request due to error when reading JSON input: %v", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		decoder := json.NewDecoder(bytes.NewBuffer(body))
		if err := decoder.Decode(&newScreening); err != nil {
			logger.Printf("ERROR rejecting as bad request due to error when parsing JSON input: %v", err)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		// make service call to create screenings
		screeningID, err := service.Create(req.Context(), NewIdent(subject), newScreening)
		if err != nil {
			logger.Printf("ERROR internal server error encountered when creating screening: %v", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		response, _ := json.Marshal(&struct {
			ScreeningID string `json:"screeningID"`
		}{screeningID.String()})
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		res.Write(response)
		logger.Printf("INFO POST %s request succeeded: %v %s", req.URL.Path, http.StatusCreated, "Created")
	}
}
