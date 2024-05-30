package requests

import (
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"

	"github.com/sirupsen/logrus"
)

// any - interface{}
// type Person struct {} - var persons []Peson
// {
//		{name surname}
//		{name surname}
// } json Marshall - preobrazuet interface to json format

type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

// UnmarshalEasyJSON implements easyjson.Unmarshaler.
func (r *Response) UnmarshalEasyJSON(w *jlexer.Lexer) {
	panic("unimplemented")
}

// func (r *Response) UnmarshalEasyJSON(data []byte) error {
// 	return easyjson.UnmarshalEasyJSON(data, r)
// }

// UnmarshalEasyJSON implements easyjson.Unmarshaler.

func SendResponse[T easyjson.Marshaler](w http.ResponseWriter, r *http.Request, code int, Body T) {
	jsonResponse, err := easyjson.Marshal(Body)
	if err != nil {
		logrus.Info(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "https://jimder.ru")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Csrft")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Access-Control-Max-Age", "86400")

	w.WriteHeader(code)
	_, err = w.Write(jsonResponse)
	if err != nil {
		logrus.Info(err.Error())
		return
	}
}

func SendSimpleResponse(w http.ResponseWriter, _ *http.Request, code int, Body string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "https://jimder.ru")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Csrft")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Access-Control-Max-Age", "86400")

	w.WriteHeader(code)
	if _, err := w.Write([]byte(Body)); err != nil {
		logrus.Info(err.Error())
		return
	}
}
