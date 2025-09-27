package a

import (
	"io"
	"net/http"

	"config"
)

var DefaultClient = &http.Client{}

func reassignPattern() {
	io.EOF = nil

	config.DefaultClient = nil // want "reassigning variable"
	DefaultClient = nil

	http.DefaultClient = nil    // want "reassigning variable"
	http.DefaultTransport = nil // want "reassigning variable"
}
