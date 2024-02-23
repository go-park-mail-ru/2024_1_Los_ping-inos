package delivery

import (
	"fmt"
	d "main.go/internal/service"
	"net/http"
)

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
