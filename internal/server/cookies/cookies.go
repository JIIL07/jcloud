package cookies

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
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
	_ = session.Save(r, w)
}

func SetSession(w http.ResponseWriter, r *http.Request, username string) {
	session, _ := Store.Get(r, "user-session")
	session.Values["user"] = true
	session.Values["username"] = username
	_ = session.Save(r, w)
}

func GetSession(r *http.Request, name string) *sessions.Session {
	session, _ := Store.Get(r, name)
	return session
}
