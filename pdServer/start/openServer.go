package start

import (
	"NOS/pdServer/search"
	"log"
	"net/http"
)

func StartPD(endPoint string, logger *log.Logger)  {
	http.HandleFunc("/search/", search.Seacher)
	http.ListenAndServe(endPoint, nil)
}
