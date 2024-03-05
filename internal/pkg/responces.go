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

	w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
    w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.WriteHeader(code)
	_, err = w.Write(jsonResponse)
	if err != nil {
		// tut doljna bit obrabotka oshibki loggerom
		return
	}
}
