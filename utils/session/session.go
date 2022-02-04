package session

import (
	"api.ethscrow/models"
	"api.ethscrow/utils"
	"context"
	"github.com/gorilla/sessions"
	"net/http"
)

var Store = sessions.NewCookieStore([]byte(utils.SESSION_KEY))

func InitSessionStore() {
	Store.Options.HttpOnly = true
	Store.Options.Secure = true
	Store.Options.SameSite = http.SameSiteNoneMode
}

func ProtectedRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session.id")
		authenticated := session.Values["authenticated"]
		if authenticated == nil || authenticated == false {
			utils.Error(w, http.StatusForbidden, "Unauthorized")
			return
		} else {
			user := &models.User{
				Username:     session.Values["username"].(string),
				EncPublicKey: session.Values["enc_public_key"].(string),
			}
			ctx := context.WithValue(r.Context(), "user", *user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
