package cookies

import (
	"github.com/gorilla/sessions"
	"net/http"
	"os"
	"time"
)

var Store *sessions.CookieStore

func SetNewCookieStore() {
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_TOKEN")))
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
		Secure:   false,
		HttpOnly: true,
	}
}

func ClearSession(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "user-session")
	session.Options.MaxAge = -1
	session.Save(r, w) // nolint:errcheck
}

func SetSession(w http.ResponseWriter, r *http.Request, username string) {
	session, _ := Store.Get(r, "user-session")
	session.Values["user"] = true
	session.Values["username"] = username
	session.Save(r, w) // nolint:errcheck
}

func GetSession(r *http.Request, name string) *sessions.Session {
	session, _ := Store.Get(r, name)
	return session
}
