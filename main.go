package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
)

var (
	basePath     = flag.String("basepath", "files", "Base path to serve files from")
	listenAddr   = flag.String("listen", ":5000", "Listen address")
	handler      *LatestHandler
	lastFilename string
)

type LatestHandler struct {
	filepath  string
	filename  string
	newestMod time.Time
	ready     bool
}

func LoggingHandler(h http.Handler) http.Handler {
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	// Log to both the file and stdout
	return handlers.LoggingHandler(io.MultiWriter(logFile, os.Stdout), h)
}
func main() {
	flag.Parse()
	handler = new(LatestHandler)
	go handler.run()
	log.Printf("Serving directory: \"%v\" from %v", *basePath, *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, LoggingHandler(handler)))
}

func checkFile(file string, fileinfo os.FileInfo, err error) error {
	if err != nil {
		log.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}
	if fileinfo.IsDir() {
		return nil // not a file.  ignore.
	}
	if fileinfo.ModTime().After(handler.newestMod) {
		handler.newestMod = fileinfo.ModTime()
		handler.filepath = file
		handler.filename = fileinfo.Name()
	}
	return nil
}

func (l *LatestHandler) run() {
	for {
		filepath.Walk(*basePath, checkFile)
		if l.filepath != "" {
			l.ready = true
		}
		if l.filepath != lastFilename {
			log.Println("Newest file is", l.filepath)
			lastFilename = l.filepath
		}
		time.Sleep(5 * time.Second)
	}
}
func (l *LatestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if l.ready {
		w.Header().Add("X-Filename", l.filename)
		http.ServeFile(w, req, l.filepath)
	} else {
		fmt.Fprintf(w, "not ready")
	}
}
