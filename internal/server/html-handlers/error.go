package handlers

import (
	"fmt"
	"net/http"
	"net/url"
)

func ServeError(w http.ResponseWriter, r *http.Request, message string, status int) {
	q := url.Values{}
	q.Add("message", message)
	q.Add("status", fmt.Sprintf("%d", status))

	http.Redirect(w, r, "/static/error.html?"+q.Encode(), http.StatusSeeOther)
}
