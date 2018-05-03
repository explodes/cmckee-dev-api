package identity

import (
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

var oauthConf = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Scopes:       []string{"user:email"},
	Endpoint:     githuboauth.Endpoint,
}

var oauthStateString = RandomString(16)

func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *App) GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	if state != oauthStateString {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		log.Printf("couldnt retrieve 'code' from callback for token exchange")
		http.Error(w, "Couldnt retrieve code", http.StatusInternalServerError)
		return
	}

	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(oauth2.NoContext, "")
	if err != nil {
		log.Printf("%v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// TODO: come up with session/value pattern

	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("store.Get() failed with '%s'\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["BearerToken"] = token.AccessToken
	session.Save(r, w)

	log.Println(*user.Email)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (app *App) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		log.Printf("store.Get() failed with '%s'\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	session.Options.MaxAge = -1
	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
