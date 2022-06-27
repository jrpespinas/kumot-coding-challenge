package rest

import (
	"encoding/json"
	"net/http"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/listing"
)

// Define an interface of an entry point to our application
type Controller interface {
	// Accept POST body, validate the request, then return appropriate response
	ShowDetails(w http.ResponseWriter, r *http.Request)
}

// Define a Controller struct that implements the handler
type controller struct {
	service listing.Service
}

// Define a constructor to inject the Service
func NewController(service listing.Service) Controller {
	return &controller{
		service: service,
	}
}

// List github username basic details
// Route: POST /usernames
// Access: public
func (c *controller) ShowDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get POST body
	var data Users
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
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
	err = c.service.Validate(data.Usernames)
	if err != nil {
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
	details, err := c.service.ShowDetails(data.Usernames)
	if err != nil {
		code := http.StatusInternalServerError
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(ErrorResponse{
			Status:  "fail",
			Code:    code,
			Message: err.Error(),
		})
	} else {
		code := http.StatusOK
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(SuccessResponse{
			Status: "success",
			Code:   code,
			Data:   details,
		})
	}
}
