package routes

import (
	"encoding/json"
	"fmt"
	"gorfe/structs"
	"gorfe/themes"
	"net/http"
)

func GenerateRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data structs.GenerateRequest
	_ = json.NewDecoder(r.Body).Decode(&data)

	var duration float64
	var file string

	if data.Theme == "grid" {
		duration, file = themes.GenerateGridImage(data)
	}

	mapIndex := map[string]interface{}{"file": file, "duration": duration}
	marshal, _ := json.Marshal(mapIndex)
	fmt.Fprintln(w, string(marshal))
}
