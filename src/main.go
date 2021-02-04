package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorfe/routes"
	"gorfe/themes"
	"gorfe/utils"
	"log"
	"net/http"
)

func main() {
	utils.SetupConfig()

	server()
}

func server() {
	config := utils.GetConfig()
	log.Println("Starting web server on port " + config.Port + "...")

	routes.InitializeMetadataRoute()

	themes.InitializeGridTheme()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", routes.IndexRoute)
	router.HandleFunc("/metadata", routes.MetadataRoute)
	router.HandleFunc("/generate", routes.GenerateRoute)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":"+config.Port, handler))
}
