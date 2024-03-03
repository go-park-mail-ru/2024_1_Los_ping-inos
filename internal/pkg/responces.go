package requests

import (
	"encoding/json"
	"net/http"
)

// any - interface{}
// type Person struct {} - var persons []Peson
// {
//		{name surname}
//		{name surname}
// } json Marshall - preobrazuet interface to json format

func SendResponse(w http.ResponseWriter, r *http.Request, code int, Body any) {
	jsonResponse, err := json.Marshal(Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonResponse)
	if err != nil {
		// tut doljna bit obrabotka oshibki loggerom
		return
	}
}
