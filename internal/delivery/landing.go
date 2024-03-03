package delivery

import (
	"net/http"
)

func (deliver *Deliver) GetCardsHandler(mux *http.ServeMux) {
	mux.HandleFunc("/",
		func(respWriter http.ResponseWriter, request *http.Request) {

			session, err := request.Cookie("session_id") // проверка авторизации
			if err != nil || session == nil || !deliver.auth.IsAuthenticated(session.Value) {
				http.Error(respWriter, "forbidden", http.StatusForbidden)
			}

			_, err = deliver.serv.GetCards(session.Value)

			if err != nil {
				http.Error(respWriter, "can't return cards", http.StatusInternalServerError)
			} else {
				// TODO вернуть карточки
			}
		})
}
