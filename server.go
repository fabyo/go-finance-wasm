package main

import (
	"log"
	"mime"
	"net/http"
)

func main() {
	// Garante que .wasm vai com o MIME correto
	mime.AddExtensionType(".wasm", "application/wasm")

	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	addr := ":9000"
	log.Println("Servidor Go rodando em http://localhost" + addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
