package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	var (
		nameFlag       string
		urlFlag        string
		listenFlag     string
		pubKeyFileFlag string
		secKeyFileFlag string
	)
	flag.StringVar(&nameFlag, "name", "minisig.me", "Server name (default: minisig.me)")
	flag.StringVar(&urlFlag, "url", "https://minisig.me", "Base URL of the server (default: https://minisig.me)")
	flag.StringVar(&listenFlag, "listen", ":8080", "Listen address of the server (default: :8080)")
	flag.StringVar(&pubKeyFileFlag, "p", "minisign.pub", "Public key file (default: minisign.pub")
	flag.StringVar(&secKeyFileFlag, "s", filepath.Join(os.Getenv("HOME"), ".minisign/minisign.key"), "Secret key file (default: $HOME/.minisign/minisign.key")

	flag.Parse()

	log.Println("Loading keys...")

	signer, err := NewSigner(urlFlag, secKeyFileFlag, pubKeyFileFlag, "")
	if err != nil {
		log.Fatal(err)
	}

	config := ServerConfig{
		Name:    nameFlag,
		BaseURL: urlFlag,
		Signer:  signer,
	}

	server := NewServer(config)

	log.Printf("Listening on %s...", listenFlag)

	handler := handlers.CombinedLoggingHandler(os.Stdout, &server)
	log.Fatal(http.ListenAndServe(listenFlag, handler))
}
