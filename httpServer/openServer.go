package httpServer

import (
	"NOS/objects"
	"log"
	"net/http"
)

func OpenServer(endPoint string, logger *log.Logger)  {
	http.HandleFunc("/", objects.Handler)
	http.ListenAndServe(endPoint, nil)
}
