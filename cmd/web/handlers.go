package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Home-Handler zum Bedienen der Startseite
func home(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob der angeforderte Pfad genau "/" ist, 
	// wenn nicht, sendet die Funktion eine 404 Not Found Antwort an den Benutzer.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Liste der Template-Dateien für die Startseite
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// Parse die Template-Dateien
	ts, err := template.ParseFiles(files...)
	// Wenn ein Fehler beim Parsen auftritt, logge den Fehler und sende eine HTTP 500 Antwort.
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Ausführen des Templates, Fehlerbehandlung ähnlich wie beim Parsen
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// snippetView-Handler zum Anzeigen eines bestimmten Snippets
func snippetView(w http.ResponseWriter, r *http.Request) {
	// Extrahieren des Snippet-ID Parameters aus dem URL Query String,
	// bei einem Fehler oder wenn die ID kleiner als 1 ist, wird eine 404 Not Found Antwort gesendet.
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Für jetzt senden wir einfach eine Platzhalterantwort zurück
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// snippetCreate-Handler zum Erstellen eines neuen Snippets
func snippetCreate(w http.ResponseWriter, r *http.Request){
	// Überprüfen ob die Anfrage die Methode POST hat, wenn nicht, senden wir einen '405 Method Not Allowed' Status 
	// und setzen den 'Allow'-Header auf 'POST'
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Für jetzt senden wir einfach eine Platzhalterantwort zurück
	w.Write([]byte("Create a new snippet..."))
}
