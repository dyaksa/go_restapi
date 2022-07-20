package helper

import (
	"encoding/json"
	"net/http"
)

func JSONDecoder(r *http.Request, result interface{}) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	PanicIf(err)
}

func JSONEncoder(w http.ResponseWriter, response interface{}) {
	w.Header().Add("content-type", "application/json")
	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	PanicIf(err)
}
