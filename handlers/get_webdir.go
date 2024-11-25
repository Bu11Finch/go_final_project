package handlers

import (
	"log"
	"net/http"
)

const webDir = "./web"

func (h *Handlers) GetWebDir() http.Handler {
	log.Printf("Загружены файлы из %s\n", webDir)

	return http.FileServer(http.Dir(webDir))
}
