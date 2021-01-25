package routes

import (
	"encoding/json"
	"gorfe/constants"
	"gorfe/utils"
	"net/http"
)

type MetadataRouteResponse struct {
	name    string   `json:"name"`
	engine  string   `json:"engine"`
	version float32  `json:"version"`
	scheme  float32  `json:"scheme"`
	themes  []string `json:"themes"`
}

var config utils.ConfigFile

func InitializeMetadataRoute() {
	config = utils.GetConfig()
}

func MetadataRoute(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"name":    config.Name,
		"engine":  constants.EngineName,
		"version": constants.EngineVersion,
		"scheme":  constants.EngineSchemeVersion,
		"themes":  constants.EngineThemes,
	}

	json.NewEncoder(w).Encode(response)
}
