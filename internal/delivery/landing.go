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
				//TODO log + http response
			}

			fmt.Fprintf(w, "cool ids:\n")
			for i := range ids {
				fmt.Fprintf(w, "%v\n", i)
			}
		})
}
