package authentication

import (
	"github.com/gorilla/securecookie"
	"github.com/labstack/echo/v4"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func SetSession(c echo.Context, userName string) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/admin/",
		}
		c.SetCookie(cookie)
	}
}

func GetUserName(c echo.Context) (userName string) {
	if cookie, err := c.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func ClearSession(c echo.Context) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/admin/",
		MaxAge: -1,
	}
	c.SetCookie(cookie)
}
