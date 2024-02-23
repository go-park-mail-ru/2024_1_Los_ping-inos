package delivery

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"main.go/config"
	d "main.go/internal/service"
)

func StartServer() error {
	mux := http.NewServeMux()

	// сюда добавлять хендлеры страничек
	mux.HandleFunc("/", landing)

	server := http.Server{
		Addr:         config.Cfg.Server.Host + config.Cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  config.Cfg.Server.Timeout * time.Second,
		WriteTimeout: config.Cfg.Server.Timeout * time.Second,
	}

	logrus.Printf("starting server at %v", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

// Функция обработки страницы
func landing(w http.ResponseWriter, _ *http.Request) {
	ids, err := d.GetCoolIdsList()
	if err != nil {
		// а вот тут хз как ошибку обрабатывать
		// просто в лог писать?
	}

	fmt.Fprintf(w, "cool ids:\n")
	for i := range ids {
		fmt.Fprintf(w, "%v\n", i)
	}
}
