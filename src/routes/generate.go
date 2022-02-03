package routes

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"gorfe/structs"
	"gorfe/themes"
	"gorfe/utils"
	"net/http"
)

func GenerateRoute(w http.ResponseWriter, r *http.Request) {
	config := utils.GetConfig()
	var data structs.GenerateRequest
	_ = json.NewDecoder(r.Body).Decode(&data)

	//sentryTraceId := r.Header.Get("sentry-trace")

	span := sentry.StartSpan(r.Context(), "generation", sentry.ContinueFromRequest(r))

	fmt.Println(span.TraceID)

	var duration float64
	var file string

	if data.Theme == "grid" {
		duration, file = themes.GenerateGridImage(data, span)
	}

	span.Finish()

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
