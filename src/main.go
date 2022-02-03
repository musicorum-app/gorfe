package main

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorfe/routes"
	"gorfe/themes"
	"gorfe/utils"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	utils.SetupConfig()

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		Debug:            true,
		TracesSampleRate: 1.0,
		Environment:      env,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	server()
}

func server() {
	config := utils.GetConfig()

	log.Println("Starting web server on port " + config.Port + "...")

	routes.InitializeMetadataRoute()

	themes.InitializeGridTheme()

	router := mux.NewRouter().StrictSlash(true)

	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	router.Handle("/", sentryHandler.Handle(http.DefaultServeMux))

	router.HandleFunc("/", routes.IndexRoute)
	router.HandleFunc("/metadata", routes.MetadataRoute)
	router.HandleFunc("/generate", routes.GenerateRoute)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	defer sentry.Flush(2 * time.Second)

	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}
