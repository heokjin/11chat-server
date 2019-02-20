package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger(NAME_PROJECT)
var hub *Hub

func main() {
	fmt.Println("...")

	runtime.GOMAXPROCS(runtime.NumCPU())

	initLog()
	hub = newHub()
	go hub.run()

	customLoggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_custom}", "remote_ip":"${remote_ip}", "host":"${host}", "method":"${method}","path":"${path}","status":"${status}",duration":"${latency_human}"` +
			`, "user_agent":${user_agent}}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05 -0700",
		Output:           os.Stdout,
	}

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(customLoggerConfig))

	log.Debug("START")
	e.Static("/", "public")
	e.GET("/ws", serveWs)
	e.Start(PORT)
}

func initLog() {
	//Code debuging을 위한 세팅
	backend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetLevel(setLogLevel("DEBUG"), NAME_PROJECT)
	logging.SetBackend(logging.NewBackendFormatter(backend, setLogFormat(LOG_FORMAT)))
}

func setLogLevel(level string) logging.Level {
	var logLevel logging.Level
	level = strings.ToUpper(level)
	switch level {
	case "CRITICAL":
		logLevel = logging.CRITICAL
	case "ERROR":
		logLevel = logging.ERROR
	case "WARNING":
		logLevel = logging.WARNING
	case "NOTICE":
		logLevel = logging.NOTICE
	case "INFO":
		logLevel = logging.INFO
	case "DEBUG":
		logLevel = logging.DEBUG
	default:
		logLevel = logging.ERROR
	}
	return logLevel
}

func setLogFormat(format string) logging.Formatter {
	return logging.MustStringFormatter(format)
}
