package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func IndexRoute(w http.ResponseWriter, _ *http.Request) {
	mapIndex := map[string]string{"working": "ok"}
	marshal, _ := json.Marshal(mapIndex)
	fmt.Fprintln(w, string(marshal))
}