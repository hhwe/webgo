package webgo

import (
	"net/http"
	"strings"
	"time"

	"github.com/hhwe/webgo/logging"
)

var Logger *logging.Logger

func init() {
	Logger = logging.NewLogger()
	Logger.SetLever("debug")
	Logger.SetColor(true)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		Logger.Info(strings.Join([]string{r.Method, r.URL.Path, time.Since(start).String()}, " | "))
	})
}
