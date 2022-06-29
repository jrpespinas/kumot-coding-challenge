package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/http/rest"
)

// Define a router interface
type Route interface {
	Router() http.Handler
}

// Implement a router struct to implement the Route interface
type routes struct {
	controller rest.Controller
}

// Define a constructor to inject the Controller layer
func NewRouter(controller rest.Controller) Route {
	return &routes{
		controller: controller,
	}
}

// Using Chi Router, define a `/users` endpoint implemented by the controller
func (route *routes) Router() http.Handler {
	r := chi.NewRouter()

	// Define our only endpoint
	r.Post("/users", route.controller.ShowDetails)
	r.Get("/generate-token", route.controller.GenerateToken)
	return r
}
