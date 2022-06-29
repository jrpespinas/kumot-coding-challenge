package rest

import (
	"encoding/json"
	"net/http"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/listing"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/session"
	"github.com/rs/zerolog"
)

// Define an interface of an entry point to our application
type Controller interface {
	// Accept POST body, validate the request, then return appropriate response
	ShowDetails(w http.ResponseWriter, r *http.Request)

	// Generate session token to access the API
	GenerateToken(w http.ResponseWriter, r *http.Request)
}

// Define a Controller struct that implements the handler
type controller struct {
	service listing.Service
	sess    session.Service
	logger  zerolog.Logger
}

// Define a constructor to inject the Service and logger
func NewController(service listing.Service, logger zerolog.Logger, sess session.Service) Controller {
	return &controller{
		service: service,
		sess:    sess,
		logger:  logger,
	}
}

// List github user basic details
// Route: GET /generate-token
// Access: protected
func (c *controller) GenerateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Generate random token
	c.logger.Info().Str("layer", "controller").Msg("generating token")
	tokenString, err := c.sess.GenerateToken()
	if err != nil {
		c.logger.Error().Str("layer", "controller").Msgf("error: %v", err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "fail",
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	// Return token to user
	c.logger.Info().Str("layer", "controller").Msg("return token to user")
	code := http.StatusOK
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(SuccessResponse{
		Status: "success",
		Code:   code,
		Data: Token{
			TokenString: tokenString,
		},
	})

}

// List github user basic details
// Route: POST /users
// Access: protected
func (c *controller) ShowDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Session-Token")

	// Check session
	c.logger.Info().Str("layer", "controller").Msg("verifying session")
	err := c.sess.VerifyToken(tokenString)
	if err != nil {
		c.logger.Error().Str("layer", "controller").Msgf("error: %v", err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "fail",
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	// Get POST body
	c.logger.Info().Str("layer", "controller").Msg("decoding post body")
	var data Users
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		c.logger.Error().Str("layer", "controller").Msgf("error: %v", err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "fail",
			Code:    code,
			Message: "unable to decode your request",
		})
		return
	}

	// Validate `data`
	c.logger.Info().Str("layer", "controller").Msg("validating post body")
	err = c.service.Validate(data.Usernames)
	if err != nil {
		c.logger.Error().Str("layer", "controller").Msgf("error: %v", err.Error())
		code := http.StatusBadRequest
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "fail",
			Code:    code,
			Message: err.Error(),
		})
		return
	}

	// Show details of each username
	c.logger.Info().Str("layer", "controller").Msg("collecting user details")
	details, err := c.service.ShowDetails(data.Usernames)
	if err != nil {
		c.logger.Error().Str("layer", "controller").Msgf("error: %v", err.Error())
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "fail",
			Code:    code,
			Message: err.Error(),
		})
	} else {
		c.logger.Info().Str("layer", "controller").Msg("displaying user details")
		code := http.StatusOK
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(SuccessResponse{
			Status: "success",
			Code:   code,
			Data:   details,
		})
	}
}
