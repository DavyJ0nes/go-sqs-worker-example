package srv

import (
	"io"
	"log"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("indexHandler")
	io.WriteString(w, `{"Message": "all gravy baby..."}`)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("apiHandler")
	io.WriteString(w, `{"Message": "hidden message"}`)
}
