package cookies

import (
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func SetNewCookieStore() {
	Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_TOKEN")))
	Store.Options = &sessions.Options{
		MaxAge:   86400, //24 hours
		Secure:   false,
		HttpOnly: true,
	}
}
