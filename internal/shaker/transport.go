package shaker

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Stats struct {
	sync.RWMutex
	Uptime       time.Time         `json:"uptime"`
	RequestCount uint64            `json:"requests"`
	SearchCounts map[string]uint64 `json:"searches"`
}

func NewStats() *Stats {
	return &Stats{
		Uptime:       time.Now(),
		SearchCounts: map[string]uint64{},
	}
}

// Process is the middleware function for stats.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if err := next(ctx); err != nil {
			ctx.Error(err)
		}
		s.Lock()
		defer s.Unlock()
		s.RequestCount++
		query := ctx.QueryParam("q")
		if len(query) > 0 {
			s.SearchCounts[query]++
		}
		return nil
	}
}

// Handle is the endpoint to return stats to the frontend.
func (s *Stats) Handle(c echo.Context) error {
	s.RLock()
	defer s.RUnlock()
	return c.JSON(http.StatusOK, s)
}

// RegisterHTTP handles routes related to shaker service
func RegisterHTTP(e *echo.Echo, svc Service) {
	e.GET("/stats", NewStats().Handle)

	e.GET("/search", func(ctx echo.Context) error {
		ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		query := ctx.QueryParam("q")
		if len(query) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "missing search query in URL params")
		}

		results := svc.Search(strings.ToLower(query))
		enc := json.NewEncoder(ctx.Response())

		ctx.Response().WriteHeader(http.StatusOK)
		// we're streaming json response, as soon as it arrives via string channel
		for result := range results {
			if err := enc.Encode(result); err != nil {
				// oh, no! headers already sent... cannot return http.StatusInternalServerError.
				// Just log the message and continue, since there isn't much we can do, just to deliver what we can (an lie to the user!!!)
				logrus.WithField("transport", "search").Errorf("error encoding json : %v", err)
				continue
			}
			ctx.Response().Flush()
		}
		return nil
	})
}
