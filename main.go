package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/halium-project/server/db"
	"github.com/halium-project/server/front"
	"github.com/halium-project/server/resource/accesstoken"
	"github.com/halium-project/server/resource/authorizationcode"
	"github.com/halium-project/server/resource/client"
	"github.com/halium-project/server/resource/user"
	"github.com/halium-project/server/saga/oauth2"
	"github.com/halium-project/server/util"
	"github.com/halium-project/server/util/endpoint"
	"github.com/halium-project/server/util/errors"
	"github.com/halium-project/server/util/permission"
	"github.com/rs/cors"
	"gitlab.com/Peltoche/yaccc"
)

func init() {
	// Force the UTC format for all the dates in order to avoid all the timezone mess.
	time.Local = time.UTC
}

const addr = ":42000"

func main() {
	ctx := context.Background()

	couchdb, err := yaccc.NewServer(util.MustGetEnv("COUCHDB_URL"), 5, time.Second)
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

	// Instantiate the controllers
	clientController := client.InitController(ctx, couchdb)
	userController := user.InitController(ctx, couchdb)
	accessTokenController := accesstoken.InitController(ctx, couchdb)
	authorizationCodeController := authorizationcode.InitController(ctx, couchdb)
	osinStorageController := oauth2.NewStorageController(clientController, authorizationCodeController, accessTokenController)
	oauth2SagaController := oauth2.InitController(ctx, couchdb, templateRenderer, userController, osinStorageController)

	router := mux.NewRouter()
	perm := permission.NewController(ctx, accessTokenController)

	// Expose the Client resource.
	clientHTTPHandler := client.NewHTTPHandler(clientController)
	router.HandleFunc("/clients", perm.Check("clients.write", clientHTTPHandler.Create)).Methods("POST")
	router.HandleFunc("/clients", perm.Check("clients.read", clientHTTPHandler.GetAll)).Methods("GET")
	router.HandleFunc("/clients/{clientID}", perm.Check("clients.read", clientHTTPHandler.Get)).Methods("GET")

	// Expose the User resource.
	userHTTPHandler := user.NewHTTPHandler(userController)
	router.HandleFunc("/users", perm.Check("users.write", userHTTPHandler.Create)).Methods("POST")
	router.HandleFunc("/users", perm.Check("users.read", userHTTPHandler.GetAll)).Methods("Get")
	router.HandleFunc("/users/{userID}", perm.Check("users.write", userHTTPHandler.Update)).Methods("PUT")
	router.HandleFunc("/users/{userID}", perm.Check("users.read", userHTTPHandler.Get)).Methods("GET")

	// Expose the OAuth2 endpoint.
	router.HandleFunc("/oauth2/token", oauth2SagaController.Token)
	router.HandleFunc("/oauth2/auth", oauth2SagaController.Authorize)
	router.HandleFunc("/oauth2/info", oauth2SagaController.Info)

	// Expose utility endpoints.
	router.HandleFunc("/ping", endpoint.Pinger).Methods("GET")

	// Link all the middlewares together and create the HTTP server.
	serv := handlers.LoggingHandler(os.Stdout, router)
	serv = cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	}).Handler(serv)

	// Start the server.
	log.Printf("listen on: %q", addr)
	err = http.ListenAndServe(addr, serv)
	if err != nil {
		log.Fatal(err)
	}
}
