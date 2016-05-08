package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/sequoiia/twiVod/server/controller"
)

func main() {
	httpCli := http.DefaultClient
	controller.HttpClient = httpCli

	n := negroni.New(negroni.NewRecovery())
	n.UseHandler(newRouter())

	n.Run("0.0.0.0:32499")
}

func newRouter() *mux.Router {
	router := mux.NewRouter()

	router.Handle("/", negroni.New(
		negroni.Wrap(http.HandlerFunc(rootHandle)),
	))

	router.PathPrefix("/api").Handler(negroni.New(
		negroni.Wrap(newApiRouter()),
	))

	return router
}

func newApiRouter() *mux.Router {
	routerBase := mux.NewRouter()

	router := routerBase.PathPrefix("/api").Subrouter()

	router.Handle("/lookup/user/{TwitchUsername}", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.LookupTwitchUser)),
	))

	router.Handle("/lookup/user/{TwitchUsername}/vods", negroni.New(
		negroni.Wrap(http.HandlerFunc(controller.LookupTwitchUserVods)),
	))

	return routerBase
}



func rootHandle(w http.ResponseWriter, r *http.Request) {
}