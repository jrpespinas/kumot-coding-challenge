package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jrpespinas/kumot-coding-challenge/pkg/http/rest"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/listing"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/logging"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/repository/api"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/repository/cache"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/router"
	"github.com/jrpespinas/kumot-coding-challenge/pkg/session"
	"github.com/rs/zerolog"
)

var (
	host       string           = fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	caching    cache.Repository = cache.NewRedis(host, 0, 2, 5)
	repository api.Repository   = api.NewRepository(os.Getenv("GITHUB_URL"))
	logger     zerolog.Logger   = logging.NewLogger(os.Getenv("LOG_LEVEL"))
	service    listing.Service  = listing.NewService(repository, caching, logger)
	sess       session.Service  = session.NewService(caching)
	controller rest.Controller  = rest.NewController(service, logger, sess)
	httpRouter router.Route     = router.NewRouter(controller)
)

func main() {
	// Set port
	logger.Info().Msgf("Redis server running at %v", host)
	logger.Info().Msgf("Serving at %v", os.Getenv("PORT"))
	port := fmt.Sprintf(":%v", os.Getenv("PORT"))

	// Run server
	http.ListenAndServe(port, httpRouter.Router())
}
