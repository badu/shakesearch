package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/oklog/oklog/pkg/group"
	"github.com/sirupsen/logrus"

	"pulley.com/shakesearch/internal/shaker"
)

func setupLogger(logLevel string) {
	switch logLevel {
	case "EMERG":
		logrus.SetLevel(logrus.ErrorLevel)
	case "ALERT":
		logrus.SetLevel(logrus.ErrorLevel)
	case "CRIT":
		logrus.SetLevel(logrus.ErrorLevel)
	case "ERR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "WARNING":
		logrus.SetLevel(logrus.WarnLevel)
	case "NOTICE":
		logrus.SetLevel(logrus.InfoLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	}
}

var (
	host     string
	port     int
	logLevel string
)

// ServerHeader middleware adds a `Server` header to the response, just for fun.
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "Shakespeare/0.1.0")
		return next(c)
	}
}

func main() {
	flag.StringVar(&host, "host", "localhost", "the host on which this server runs (default : 'localhost')")
	flag.IntVar(&port, "port", 8080, "the port to listen on for insecure connections, defaults to a 8080")
	flag.StringVar(&logLevel, "loglevel", "DEBUG", "the log level for the logrus logger, defaults to ERR. Valid values are : 'EMERG','ALERT','CRIT','ERR','WARNING','NOTICE','INFO' and 'DEBUG'")

	srvFields := logrus.Fields{"service": "shakespeare"}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	dir, err := os.Getwd()
	if err != nil {
		logrus.WithFields(srvFields).Errorf("error getting current folder : %v", err)
		os.Exit(1)
	}
	dataFileName := dir + "/data/completeworks.txt"

	logrus.WithFields(srvFields).Debugf("data file was found in %s", dataFileName)

	// yeah, I know echo comes with a logger, but why not having two of them :)
	setupLogger(logLevel)

	// I like this group, since two years ago
	var g group.Group
	// wiring things up : I like a cup of good composition in the morning
	{
		repository, err := shaker.NewRepository(dataFileName)
		if err != nil {
			logrus.WithFields(srvFields).Errorf("error creating repository : %v", err)
			os.Exit(1)
		}
		var (
			service = shaker.NewService(repository)
			router  = echo.New()
		)
		// Middleware
		router.Use(middleware.Recover())
		router.Use(middleware.CORS())
		router.Use(ServerHeader)
		router.Use(middleware.Static(dir + "/web/dist"))

		shaker.RegisterHTTP(router, service)

		g.Add(func() error {
			logrus.WithFields(srvFields).Debugf("HTTP listening at %s:%d", host, port)
			return router.Start(fmt.Sprintf("%s:%d", host, port))
		}, func(error) {
			_ = router.Close()
		})
	}

	// This function just sits and waits for kill switch, like a good cloud citizen
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("signal %s received", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	_ = g.Run()
}
