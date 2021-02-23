package main

import (
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/database"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/flags"
	"github.com/kabukky/journey/logger"
	"github.com/kabukky/journey/plugins"
	"github.com/kabukky/journey/repositories/file"
	"github.com/kabukky/journey/server"
	"github.com/kabukky/journey/structure/methods"
	"github.com/kabukky/journey/templates"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.HideBanner = true
	logger.DefaultLogger = e.Logger

	_ = initComponents()
	loadRuteConfig(e)

	// HTTP(S) Server
	httpPort := configuration.Config.HttpHostAndPort
	httpsPort := configuration.Config.HttpsHostAndPort
	// Check if HTTP/HTTPS flags were provided
	if flags.HttpPort != "" {
		components := strings.SplitAfterN(httpPort, ":", 2)
		httpPort = components[0] + flags.HttpPort
	}
	if flags.HttpsPort != "" {
		components := strings.SplitAfterN(httpsPort, ":", 2)
		httpsPort = components[0] + flags.HttpsPort
	}

	logger.Info("Starting server without HTTPS support. Please enable HTTPS in " + filenames.ConfigFilename + " to improve security.")

	logger.Info("Starting https server on port " + httpsPort + "...")
	go func() {
		if err := e.StartAutoTLS(httpsPort); err != nil {
			logger.Fatal("Error: Couldn't start the HTTPS server:", err)
		}
	}()

	logger.Info("Starting http server on port " + httpPort + "...")
	if err := e.Start(httpPort); err != nil {
		return
	}
}

func initComponents() (err error) {
	if err = database.Initialize(); err != nil {
		logger.Fatal("Error: Couldn't initialize database:", err)
		return
	}

	if err = methods.GenerateBlog(); err != nil {
		logger.Fatal("Error: Couldn't generate blog data:", err)
		return
	}

	if err = templates.Generate(); err != nil {
		logger.Fatal("Error: Couldn't compile templates:", err)
		return
	}

	if err = plugins.Load(); err == nil {
		// Close LuaPool at the end
		defer plugins.LuaPool.Shutdown()
		logger.Info("Plugins loaded.")
	}

	if err = file.InitFS(); err != nil {
		logger.Fatal("Error: Couldn't init filesystem:", err)
		return
	}
	return
}

func loadRuteConfig(e *echo.Echo) {
	server.InitializeBlog(e)
	server.InitializePages(e)
	server.InitializeAdmin(e)
}
