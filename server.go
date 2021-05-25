package main

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

//go:embed templates/*
var content embed.FS

type ServerConfig struct {
	Name    string
	BaseURL string
	Signer  Signer
}

type Server struct {
	config    ServerConfig
	router    *httprouter.Router
	templates *template.Template
}

func NewServer(config ServerConfig) Server {
	router := httprouter.New()

	server := Server{
		config:    config,
		router:    router,
		templates: template.Must(template.ParseFS(content, "templates/*.tmpl")),
	}

	router.GET("/", server.indexHandler)
	router.POST("/sign", server.signHandler)
	router.GET("/minisign.pub", server.keyHandler)

	return server
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func (server *Server) indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := server.templates.ExecuteTemplate(w, "index.html.tmpl", server.config); err != nil {
		log.Println(err)
	}
}

func (server *Server) signHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	hexDigest := r.PostFormValue("digest")
	now := time.Now().UTC()

	if hexDigest == "" {
		err := errors.New("digest is empty")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	digest, err := decodeHexDigest([]byte(hexDigest))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signature, err := server.config.Signer.SignDigest(digest, now)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	w.Write(signature)
	io.WriteString(w, "\n")
}

func (server *Server) keyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := server.config.Signer.PublicKey.String()

	untrustedComment := fmt.Sprintf("%s public key", server.config.Signer.By)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	io.WriteString(w, fmt.Sprintf("untrusted comment: %s\n", untrustedComment))
	io.WriteString(w, key)
	io.WriteString(w, "\n")
}
