package server

import "net/http"

type Server struct {
	httpServer *http.Server
}
