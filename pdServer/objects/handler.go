package objects

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println(r.Method)
}
