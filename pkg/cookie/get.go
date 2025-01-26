package cookie

import (
	"errors"
	"net/http"
)

func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", errors.New("cookie not found: " + name)
		}
		return "", err
	}
	return cookie.Value, nil
}
