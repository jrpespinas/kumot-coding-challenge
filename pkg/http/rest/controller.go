package rest

import (
	"encoding/json"
	"net/http"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/listing"
	"github.com/rs/zerolog"
)

// Define an interface of an entry point to our application
type Controller interface {
	// Accept POST body, validate the request, then return appropriate response
	ShowDetails(w http.ResponseWriter, r *http.Request)
}

// Define a Controller struct that implements the handler
type controller struct {
	service listing.Service
	logger  zerolog.Logger
}

// Define a constructor to inject the Service and logger
func NewController(service listing.Service, logger zerolog.Logger) Controller {
	return &controller{
		service: service,
		logger:  logger,
	}
}

// List github user basic details
// Route: POST /users
// Access: public
func (c *controller) ShowDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get POST body
	c.logger.Info().Str("layer", "controller").Msg("decoding post body")
	var data Users
	err := json.NewDecoder(r.Body).Decode(&data)
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
