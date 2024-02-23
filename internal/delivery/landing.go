package delivery

import (
	"fmt"
	"net/http"
)

func (c *Deliver) landingHandler(mux *http.ServeMux) {
	mux.HandleFunc("/",
		func(w http.ResponseWriter, _ *http.Request) {
			ids, err := c.serv.GetCoolIdsList()
			if err != nil {
				// а вот тут хз как ошибку обрабатывать
				// просто в лог писать?
			}

			fmt.Fprintf(w, "cool ids:\n")
			for i := range ids {
				fmt.Fprintf(w, "%v\n", i)
			}
		})
}
