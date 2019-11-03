package start

import (
	"NOS/httpServer/operation"
	"log"
	"net/http"
)

func OpenServer(endPoint string, logger *log.Logger)  {
	http.HandleFunc("/", operation.Handler)
	http.ListenAndServe(endPoint, nil)
}
