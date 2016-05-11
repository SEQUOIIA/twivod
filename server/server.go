package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/sequoiia/twiVod/server/controller"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"crypto/rsa"
	"crypto/rand"
	"log"
	"fmt"
	"time"
)

var jwtMiddleware *jwtmiddleware.JWTMiddleware

var PubKey rsa.PublicKey
var PrivKey *rsa.PrivateKey

func main() {
	httpCli := http.DefaultClient
	controller.HttpClient = httpCli

	privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	PubKey = privKey.PublicKey
	PrivKey = privKey

	jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func (token *jwt.Token) (interface{}, error) {
			return &PubKey, nil
		},
		SigningMethod: jwt.SigningMethodRS512,
		Debug: false,
		Extractor: jwtmiddleware.FromCookie,
	})

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
		negroni.HandlerFunc(jwtMiddleware.HandlerWithNext),
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
	jwtToken := jwt.New(jwt.SigningMethodRS512)

	jwtToken.Claims["AccessToken"] = "level1"
	jwtToken.Claims["UserInfo"] = struct {
		Name string
		Access bool
	}{Name: "sequoiia", Access: true}

	jwtToken.Claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	jwtTokenString, err := jwtToken.SignedString(PrivKey)
	if err != nil {
		log.Fatal(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name: "AccessToken",
		Value: jwtTokenString,
		Path: "/",
		RawExpires: "0",
	})

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cookie set."))
}