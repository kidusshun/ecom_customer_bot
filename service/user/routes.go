package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/kidusshun/ecom_bot/service/auth"
	"github.com/markbates/goth/gothic"
)

type Handler struct {
	store UserStore
}

func NewHandler(store UserStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/auth/login", h.handleLogin)
	router.Get("/auth/{provider}/callback", h.getAuthCallbackFunction)
	router.Get("/logout/{provider}", h.handleLogout)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	url := auth.GoogleOauthConfig.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}
func (h *Handler) handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify state token (you can add more checks for security)
	if r.FormValue("state") != "randomstate" {
		http.Error(w, "State token mismatch", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	token, err := auth.GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve user information from Google
	client := auth.GoogleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse and return user information
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func (h *Handler) getAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(user)

	http.Redirect(w, r, "http://localhost:5173", http.StatusFound)

}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)

}

func (h *Handler) handleAuth(w http.ResponseWriter, r *http.Request) {
	// try to get the user without re-authenticating
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(w, gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}

}

var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
