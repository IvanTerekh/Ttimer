package auth

import (
	"net/http"
	"net/url"
	"os"
	"ttimer/app"
)

// LogoutHandler handles user's logout.
func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		handleError(err, w)
		return
	}
	session.Options.MaxAge = -1
	app.Store.Save(r, w, session)

	domain := os.Getenv("AUTH0_DOMAIN")

	var Url *url.URL
	Url, err = url.Parse("https://" + domain)
	if err != nil {
		handleError(err, w)
		return
	}

	Url.Path += "/v2/logout"
	parameters := url.Values{}
	parameters.Add("returnTo", "http://"+os.Getenv("TTIMER_DOMAIN"))
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	Url.RawQuery = parameters.Encode()

	http.Redirect(w, r, Url.String(), http.StatusTemporaryRedirect)
}
