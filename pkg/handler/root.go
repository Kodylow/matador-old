// handler/root_handler.go

package handler

import (
	"github.com/russross/blackfriday/v2"
	"io/ioutil"
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	md, err := ioutil.ReadFile("README.md")
	if err != nil {
		http.Error(w, "Couldn't read file", http.StatusInternalServerError)
		return
	}

	html := blackfriday.Run(md)

	w.Write(html)
}
