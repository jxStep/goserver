package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	// Definieren der Standard-Adresse für den HTTP-Server
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Erzeugen eines neuen HTTP Request Multiplexers
	mux := http.NewServeMux()

	// Erzeugen eines neuen Datei-Servers, der statische Dateien aus dem Ordner "./static" bereitstellt
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	// Statische Dateien sind nicht direkt unter "/static" zugänglich
	mux.Handle("/static", http.NotFoundHandler())
	// Statische Dateien sind unter "/static/" zugänglich
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Festlegen der Handlerfunktionen für die jeweiligen Pfade
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// Starten des Servers und Ausgeben einer Fehlermeldung, wenn der Server nicht gestartet werden kann
	log.Printf("Starting Server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}

// Definieren einer neuen Dateisystemstruktur, die das Standard-HTTP-Dateisystem umschließt
type neuteredFileSystem struct {
	fs http.FileSystem
}

// Definieren einer eigenen Open-Funktion für das neuteredFileSystem
func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	// Öffnen der angeforderten Datei
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	// Abrufen der Dateiinformationen
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Überprüfen, ob die angeforderte Datei ein Verzeichnis ist
	if s.IsDir() {
		// Wenn ja, versuchen, eine Datei namens "index.html" in diesem Verzeichnis zu öffnen
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			// Schließen der ursprünglich geöffneten Datei und Rückgabe des Fehlers
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	// Rückgabe der geöffneten Datei, wenn keine Fehler aufgetreten sind
	return f, nil
}
