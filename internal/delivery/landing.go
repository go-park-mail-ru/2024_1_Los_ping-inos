package delivery

import (
	"main.go/internal/types"
	"net/http"
	"strconv"
)

func (deliver *Deliver) GetCardsHandler(mux *http.ServeMux) {
	mux.HandleFunc("/",
		func(respWriter http.ResponseWriter, request *http.Request) {

			session, err := request.Cookie("session_id") // проверка авторизации
			if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value) {
				http.Error(respWriter, "forbidden", http.StatusForbidden)
				return
			}

			var lastID int
			if request.URL.Query().Get("last") != "" {
				lastID, err = strconv.Atoi(request.URL.Query().Get("last"))
			} else {
				lastID = 0
			}

			if err != nil {
				http.Error(respWriter, "can't return cards: no last ID", http.StatusInternalServerError)
				return
			}

			cards, err := deliver.serv.GetCards(session.Value, types.UserID(lastID))
			if err != nil {
				http.Error(respWriter, "can't return cards: smth went wrong", http.StatusInternalServerError)
				return
			}

			respWriter.Header().Set("Content-Type", "application/json")
			_, err = respWriter.Write([]byte(cards))
			if err != nil {
				http.Error(respWriter, "can't return cards: response writer error", http.StatusInternalServerError)
				return
			}
		})
}
