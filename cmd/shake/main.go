package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
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

// pprofIndex will pass the call from /debug/pprof to pprof.
func pprofIndex() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Index(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofHeap will pass the call from /debug/pprof/heap to pprof.
func pprofHeap() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("heap").ServeHTTP(ctx.Response(), ctx.Request())
		return nil
	}
}

// pprofGoroutine will pass the call from /debug/pprof/goroutine to pprof.
func pprofGoroutine() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofBlock will pass the call from /debug/pprof/block to pprof.
func pprofBlock() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("block").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofThreadCreate will pass the call from /debug/pprof/threadcreate to pprof.
func pprofThreadCreate() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofCmdline will pass the call from /debug/pprof/cmdline to pprof.
func pprofCmdline() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Cmdline(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofProfile will pass the call from /debug/pprof/profile to pprof.
func pprofProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Profile(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofSymbol will pass the call from /debug/pprof/symbol to pprof.
func pprofSymbol() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Symbol(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofTrace will pass the call from /debug/pprof/trace to pprof.
func pprofTrace() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Trace(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofMutex will pass the call from /debug/pprof/mutex to pprof.
func pprofMutex() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("mutex").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}

// pprofAllocs will pass the call from /debug/pprof/allocs to pprof.
func pprofAllocs() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pprof.Handler("allocs").ServeHTTP(ctx.Response().Writer, ctx.Request())
		return nil
	}
}
func mountPProf(router *echo.Echo) {
	routers := []struct {
		Method  string
		Path    string
		Handler echo.HandlerFunc
	}{
		{http.MethodGet, "/debug/pprof/", pprofIndex()},
		{http.MethodGet, "/debug/pprof/allocs", pprofAllocs()},
		{http.MethodGet, "/debug/pprof/heap", pprofHeap()},
		{http.MethodGet, "/debug/pprof/goroutine", pprofGoroutine()},
		{http.MethodGet, "/debug/pprof/block", pprofBlock()},
		{http.MethodGet, "/debug/pprof/threadcreate", pprofThreadCreate()},
		{http.MethodGet, "/debug/pprof/cmdline", pprofCmdline()},
		{http.MethodGet, "/debug/pprof/profile", pprofProfile()},
		{http.MethodGet, "/debug/pprof/symbol", pprofSymbol()},
		{http.MethodPost, "/debug/pprof/symbol", pprofSymbol()},
		{http.MethodGet, "/debug/pprof/trace", pprofTrace()},
		{http.MethodGet, "/debug/pprof/mutex", pprofMutex()},
	}

	for _, r := range routers {
		switch r.Method {
		case http.MethodGet:
			router.GET(r.Path, r.Handler)
		case http.MethodPost:
			router.POST(r.Path, r.Handler)
		}
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
	sqlFileName := dir + "/data/Shakespeare.sql"

	logrus.WithFields(srvFields).Debugf("data file was found in %s", dataFileName)

	// yeah, I know echo comes with a logger, but why not having two of them :)
	setupLogger(logLevel)

	// I like this group, since two years ago
	var g group.Group
	// wiring things up : I like a cup of good composition in the morning
	{
		repository, err := shaker.NewRepository(dataFileName, sqlFileName)
		if err != nil {
			logrus.WithFields(srvFields).Errorf("error creating repository : %v", err)
			os.Exit(1)
		}
		var (
			service = shaker.NewService(repository)
			router  = echo.New()
		)
		mountPProf(router)
		echo.NotFoundHandler = func(c echo.Context) error {
			// render your 404 page
			return c.String(http.StatusNotFound, "not found page")
		}
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
