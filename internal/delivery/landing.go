package delivery

import (
	"fmt"
	"net/http"
)

func (c *Deliver) GetCardsHandler(mux *http.ServeMux) {
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			_, err := c.serv.GetCards(w, r)
			if err != nil {
				http.Error(w, "forbidden", http.StatusForbidden)
				fmt.Println("access denied!")
			} else {
				// TODO вернуть карточки
				fmt.Println("u'r good")
			}
		})
}
