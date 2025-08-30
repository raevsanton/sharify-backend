package cookie

import (
	"net/http"
)

func SetCookie(w http.ResponseWriter, name, value string, maxAge int, isProduction bool) {
	var sameSite http.SameSite
	if isProduction {
		sameSite = http.SameSiteNoneMode
	} else {
		sameSite = http.SameSiteLaxMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		SameSite: sameSite,
		MaxAge:   maxAge,
		Path:     "/",
	})
}
