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

	if data.Theme == "grid" {
		themes.GenerateGridImage(data)
	}

	mapIndex := map[string]string{"working": "ok"}
	marshal, _ := json.Marshal(mapIndex)
	fmt.Fprintln(w, string(marshal))
}