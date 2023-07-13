package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Definieren der Standard-Adresse für den HTTP-Server
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

    // Use log.New() to create a logger for writing information messages. This takes
    // three parameters: the destination to write the logs to (os.Stdout), a string
    // prefix for message (INFO followed by a tab), and flags to indicate what
    // additional information to include (local date and time). Note that the flags
    // are joined using the bitwise OR operator |.
    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

    // Create a logger for writing error messages in the same way, but use stderr as
    // the destination and use the log.Lshortfile flag to include the relevant
    // file name and line number.
    errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile)

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
	}

	// Erzeugen eines neuen HTTP Request Multiplexers
	mux := http.NewServeMux()

	// Erzeugen eines neuen Datei-Servers, der statische Dateien aus dem Ordner "./static" bereitstellt
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static")})
	// Statische Dateien sind nicht direkt unter "/static" zugänglich
	mux.Handle("/static", http.NotFoundHandler())
	// Statische Dateien sind unter "/static/" zugänglich
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Festlegen der Handlerfunktionen für die jeweiligen Pfade
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Starten des Servers und Ausgeben einer Fehlermeldung, wenn der Server nicht gestartet werden kann
	infoLog.Printf("Starting Server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
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
