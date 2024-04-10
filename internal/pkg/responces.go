package requests

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
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
		logrus.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://185.241.192.216:8081")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Csrft")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Max-Age", "86400")

	w.WriteHeader(code)
	_, err = w.Write(jsonResponse)
	if err != nil {
		logrus.Info(err.Error())
		return
	}
}
