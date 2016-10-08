package negronigomol

import (
	"net/http"
	"time"

	"github.com/aphistic/gomol"
	"github.com/urfave/negroni"
)

// Logger is a Negroni middleware for logging requests with the gomol logging library
type Logger struct {
	base *gomol.Base

	LogStart bool
	LogEnd   bool

	LogLevel gomol.LogLevel

	FieldMethod string
	FieldPath   string
	FieldStatus string
	FieldTimeMs string
}

// NewLogger will create a new Logger using the current default gomol logger base
func NewLogger() *Logger {
	return &Logger{
		base: gomol.Default(),

		LogStart: true,
		LogEnd:   true,

		LogLevel: gomol.LevelInfo,

		FieldMethod: "method",
		FieldPath:   "path",
		FieldStatus: "status",
		FieldTimeMs: "time_ms",
	}
}

// NewLoggerForBase will create a new Logger using the provided gomol logger base
func NewLoggerForBase(base *gomol.Base) *Logger {
	logger := NewLogger()
	logger.base = base
	return logger
}

func (l *Logger) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	reqAttrs := gomol.NewAttrs()
	if len(l.FieldMethod) > 0 {
		reqAttrs.SetAttr(l.FieldMethod, req.Method)
	}
	if len(l.FieldPath) > 0 {
		reqAttrs.SetAttr(l.FieldPath, req.URL.Path)
	}
	start := time.Now()
	l.base.Log(l.LogLevel, reqAttrs, "Started %v %v", req.Method, req.URL.Path)

	next(rw, req)

	res := rw.(negroni.ResponseWriter)

	resAttrs := gomol.NewAttrsFromMap(reqAttrs.Attrs())
	if len(l.FieldStatus) > 0 {
		resAttrs.SetAttr(l.FieldStatus, res.Status())
	}
	if len(l.FieldTimeMs) > 0 {
		duration := time.Since(start)
		resAttrs.SetAttr(l.FieldTimeMs, float64(duration.Nanoseconds())/float64(1000000))
	}
	l.base.Log(l.LogLevel, resAttrs, "Completed %v %v (%v)", req.Method, req.URL.Path, res.Status())
}
