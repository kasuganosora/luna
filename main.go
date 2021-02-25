package main

import (
	"github.com/kabukky/journey/configuration"
	"github.com/kabukky/journey/dao"
	"github.com/kabukky/journey/dao/scheme"
	"github.com/kabukky/journey/filenames"
	"github.com/kabukky/journey/flags"
	"github.com/kabukky/journey/logger"
	"github.com/kabukky/journey/repositories/file"
	"github.com/kabukky/journey/repositories/setting"
	"github.com/kabukky/journey/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"
)

func main() {

	var err error
	setting.LoadEnv()
	var dsn *scheme.Setting
	if dsn, err = setting.GetGlobal("dsn"); err != nil {
		panic("Get DB DSN error: " + err.Error())
	}

	dao.InitDao(dsn.GetString(), setting.IsDebugMode())
	err = setting.LoadCache(dao.DB)
	if err != nil {
		logger.Fatal(err)
		return
	}

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

	//if err = plugins.Load(); err == nil {
	//	// Close LuaPool at the end
	//	defer plugins.LuaPool.Shutdown()
	//	logger.Info("Plugins loaded.")
	//}

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
