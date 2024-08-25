package server

import "net/http"

func Middleware(next http.Handler) http.Handler {
	return nil
}
