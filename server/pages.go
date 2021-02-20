package server

import (
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/helpers"
	"github.com/labstack/echo/v4"
	"net/http"
	"path/filepath"
	"strings"
)

func pagesHandler(c echo.Context) (err error) {
	path := filepath.Join(filenames.PagesFilepath, c.Param("filepath"))
	// If the path points to a directory, add a trailing slash to the path (needed if the page loads relative assets).
	if helpers.IsDirectory(path) && !strings.HasSuffix(c.Request().RequestURI, "/") {
		err = c.Redirect(http.StatusMovedPermanently, c.Request().RequestURI+"/")
		return
	}
	err = c.File(path)

	return
}

func InitializePages(router *echo.Echo) {
	// For serving standalone projects or pages saved in in content/pages
	router.GET("/pages/:filepath", pagesHandler)
}
