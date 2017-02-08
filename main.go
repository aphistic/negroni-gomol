package negronigomol

import (
	"net/http"
	"time"

	"github.com/aphistic/gomol"
	requestid "github.com/aphistic/negroni-requestid"
	"github.com/urfave/negroni"
)

// Logger is a Negroni middleware for logging requests with the gomol logging library
type Logger struct {
	base *gomol.Base

	LogStart bool
	LogEnd   bool

	LogLevel gomol.LogLevel

	FieldMethod    string
	FieldPath      string
	FieldStatus    string
	FieldTimeMs    string
	FieldRequestID string
}

// NewLogger will create a new Logger using the current default gomol logger base
func NewLogger() *Logger {
	return &Logger{
		base: gomol.Default(),

		LogStart: true,
		LogEnd:   true,

		LogLevel: gomol.LevelInfo,

		FieldMethod:    "request_method",
		FieldPath:      "path",
		FieldStatus:    "status",
		FieldTimeMs:    "time_ms",
		FieldRequestID: "request_id",
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
	if l.FieldMethod != "" {
		reqAttrs.SetAttr(l.FieldMethod, req.Method)
	}
	if l.FieldPath != "" {
		reqAttrs.SetAttr(l.FieldPath, req.URL.Path)
	}
	if l.FieldRequestID != "" {
		if reqID, err := requestid.FromContext(req.Context()); err == nil {
			reqAttrs.SetAttr(l.FieldRequestID, reqID)
		}
	}
	start := time.Now()

	// stash the path since http.StripPath will make it very confusing later
	path := req.URL.Path

	l.base.Log(l.LogLevel, reqAttrs, "Started %v %v", req.Method, path)

	next(rw, req)

	res := rw.(negroni.ResponseWriter)

	resAttrs := gomol.NewAttrsFromMap(reqAttrs.Attrs())
	if l.FieldStatus != "" {
		resAttrs.SetAttr(l.FieldStatus, res.Status())
	}
	if l.FieldTimeMs != "" {
		duration := time.Since(start)
		resAttrs.SetAttr(l.FieldTimeMs, float64(duration.Nanoseconds())/float64(1000000))
	}
	l.base.Log(l.LogLevel, resAttrs, "Completed %v %v (%v)", req.Method, path, res.Status())
}
