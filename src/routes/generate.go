package routes

import (
	"encoding/json"
	"fmt"
	"gorfe/structs"
	"gorfe/themes"
	"gorfe/utils"
	"net/http"
)

func GenerateRoute(w http.ResponseWriter, r *http.Request) {
	config := utils.GetConfig()
	var data structs.GenerateRequest
	_ = json.NewDecoder(r.Body).Decode(&data)

	var duration float64
	var file string

	if data.Theme == "grid" {
		duration, file = themes.GenerateGridImage(data)
	}

	if data.ReturnImage {
		w.Header().Set("Content-Type", "image/webp")
		http.ServeFile(w, r, config.ExportPath+file)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	mapIndex := map[string]interface{}{"file": file, "duration": duration}
	marshal, _ := json.Marshal(mapIndex)
	fmt.Fprintln(w, string(marshal))
}
