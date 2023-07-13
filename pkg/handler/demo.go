package handler

import (
	"html/template"
	"log"
	"net/http"
)

func DemoIndexHandler(w http.ResponseWriter, r *http.Request) {
	// Gather any necessary data
	log.Println("DemoIndexHandler started")
	tmpl, err := template.ParseFiles("static/main.html", "static/chat.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "main.html", map[string]string{
		"CurrentPage": "chat",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DemoChatHandler(w http.ResponseWriter, r *http.Request) {
	// Gather any necessary data
	log.Println("DemoChatHandler started")
	tmpl, err := template.ParseFiles("static/main.html", "static/chat.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "main.html", map[string]string{
		"CurrentPage": "chat",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DemoImageHandler(w http.ResponseWriter, r *http.Request) {
	// Gather any necessary data
	log.Println("DemoImageHandler started")
	tmpl, err := template.ParseFiles("static/main.html", "static/images.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "main.html", map[string]string{
		"CurrentPage": "images",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DemoEmbeddingsHandler(w http.ResponseWriter, r *http.Request) {
	// Gather any necessary data
	log.Println("DemoEmbeddingsHandler started")
	tmpl, err := template.ParseFiles("static/main.html", "static/embeddings.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "main.html", map[string]string{
		"CurrentPage": "embeddings",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
