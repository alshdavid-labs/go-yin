package yin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

var (
	green            = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white            = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow           = string([]byte{27, 91, 57, 48, 59, 52, 51, 109})
	red              = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue             = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta          = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan             = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset            = string([]byte{27, 91, 48, 109})
	consoleColorMode = autoColor
)

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	IsTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.IsTerm)
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	return fmt.Sprintf("%s %3d %s %s %-7s %s %s\n%s",
		statusColor, param.StatusCode, resetColor,
		// param.Latency,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}

// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	consoleColorMode = disableColor
}

// ForceConsoleColor force color output in the console.
func ForceConsoleColor() {
	consoleColorMode = forceColor
}

func Create(params LogFormatterParams) {
	fmt.Fprint(os.Stdout, defaultLogFormatter(params))
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := statusRecorder{w, 200}
		next.ServeHTTP(&ww, r)
		latency := time.Now().Sub(start)
		Create(LogFormatterParams{
			Path:       r.URL.String(),
			StatusCode: ww.status,
			Method:     r.Method,
			Latency:    latency,
			IsTerm:     true,
		})
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}
