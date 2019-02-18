package main

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/mux"
	"github.com/halium-project/go-server-utils/endpoint"
	"github.com/halium-project/go-server-utils/env"
	"github.com/halium-project/go-server-utils/errors"
	"github.com/halium-project/go-server-utils/server"
	"github.com/halium-project/server/db"
	"github.com/halium-project/server/front"
	"github.com/halium-project/server/resource/accesstoken"
	"github.com/halium-project/server/resource/authorizationcode"
	"github.com/halium-project/server/resource/client"
	"github.com/halium-project/server/resource/user"
	"github.com/halium-project/server/saga/oauth2"
	"github.com/halium-project/server/utils/permission"
	"gitlab.com/Peltoche/yaccc"
)

func init() {
	// Force the UTC format for all the dates in order to avoid all the timezone mess.
	time.Local = time.UTC
}

const addr = ":42000"

func main() {
	ctx := context.Background()

	couchdb, err := yaccc.NewServer(env.MustGetEnv("COUCHDB_URL"), 5, time.Second)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to connect to couchdb server"))
	}

	err = db.InitCouchdbServer(ctx, couchdb)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to init the couchdb server"))
	}

	templateRenderer, err := front.NewHTMLRenderer()
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()

	// Set the permission handler.
	accessTokenController := accesstoken.InitController(ctx, couchdb)
	perm := permission.NewController(ctx, accessTokenController)

	// Expose the Client resource.
	clientController := client.InitController(ctx, couchdb)
	clientHTTPHandler := client.NewHTTPHandler(clientController)
	router.HandleFunc("/clients", perm.Check("clients.write", clientHTTPHandler.Create)).Methods("POST")
	router.HandleFunc("/clients", perm.Check("clients.read", clientHTTPHandler.GetAll)).Methods("GET")
	router.HandleFunc("/clients/{clientID}", perm.Check("clients.read", clientHTTPHandler.Get)).Methods("GET")

	// Expose the User resource.
	userController := user.InitController(ctx, couchdb)
	userHTTPHandler := user.NewHTTPHandler(userController)
	userHTTPHandler.Register(router, perm)

	// Expose the OAuth2 endpoint.
	authorizationCodeController := authorizationcode.InitController(ctx, couchdb)
	osinStorageController := oauth2.NewStorageController(clientController, authorizationCodeController, accessTokenController)
	oauth2SagaController := oauth2.InitController(ctx, couchdb, templateRenderer, userController, osinStorageController)
	router.HandleFunc("/oauth2/token", oauth2SagaController.Token)
	router.HandleFunc("/oauth2/auth", oauth2SagaController.Authorize)
	router.HandleFunc("/oauth2/info", oauth2SagaController.Info)

	// Expose utility endpoints.
	router.HandleFunc("/ping", endpoint.Pinger).Methods("GET")

	server.ServeHandler(addr, router)
}
