package logs3

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Attr struct {
	k string
	v any
}

func Any(key string, value any) Attr {
	return Attr{k: key, v: value}
}

type Handler interface {
	// Handle processes the Record.
	// Handle methods that produce output should observe the following rules:
	//   - If r.Time() is the zero time, do not output it.
	//   - If r.Level() is Level(0), do not output it.
	Handle(Record) error

	// With returns a new Handler whose attributes consist of
	// the receiver's attributes concatenated with the arguments.
	With(attrs []Attr) Handler
}

type HandlerOptions struct {
	// If set, ReplaceAttr is called on each attribute of the message,
	// and the returned value is used instead of the original. If the returned
	// key is empty, the attribute is omitted from the output.
	//
	// The built-in attributes with keys "time", "level", "source", and "msg"
	// are passed to this function first, except that time and level are omitted
	// if zero, and source is omitted if AddSource is false.
	ReplaceAttr func(a Attr) Attr
}

type Level int

const (
	ErrorLevel Level = 10
	WarnLevel  Level = 20
	InfoLevel  Level = 30
	DebugLevel Level = 31
)

type Logger struct {
	// Has unexported fields.
	handler Handler
	ctx     context.Context
}

var defaultLogger = &Logger{
	handler: NewOtelHandler(os.Stdout),
}

func Default() *Logger {
	return defaultLogger
}

type ctxKey struct{}

func NewContext(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(ctxKey{}).(*Logger)
	if !ok {
		logger = Default()
	}
	return Logger{handler: logger.handler, ctx: ctx}
}

func New(h Handler) *Logger {
	return &Logger{handler: h}
}

func (l *Logger) Handler() Handler {
	return l.handler
}

func (l *Logger) LogAttrs(level Level, msg string, attrs ...Attr) {
	r := NewRecord(time.Now(), level, msg, 1, l.ctx)
	r.attrs = attrs
	l.handler.Handle(r)
}

// With returns a new Logger whose handler's attributes are a concatenation of
// l's attributes and the given arguments, converted to Attrs as in Logger.Log.
func (l *Logger) With(attrs ...Attr) *Logger {
	handler := l.handler.With(attrs)
	return &Logger{handler: handler}
}

type Record struct {
	t     time.Time
	level Level
	msg   string
	attrs []Attr
	ctx   context.Context
}

func NewRecord(t time.Time, level Level, msg string, calldepth int, ctx context.Context) Record {
	return Record{t: t, msg: msg, ctx: ctx}
}

type OtelHandler struct {
	w     io.Writer
	attrs map[string]any
}

func NewOtelHandler(w io.Writer) *OtelHandler {
	return &OtelHandler{w: w}
}

func (h *OtelHandler) Handle(r Record) error {
	_, err := h.w.Write([]byte(r.msg))

	spanContext := trace.SpanContextFromContext(r.ctx)

	if spanContext.IsValid() {
		_, err = h.w.Write(
			[]byte(fmt.Sprintf(
				" trace_id=%v span_id=%v",
				spanContext.TraceID(), spanContext.SpanID(),
			)),
		)
	}

	for k, v := range h.attrs {
		h.w.Write([]byte(fmt.Sprintf(" %s=%v", k, v)))
	}
	for _, attr := range r.attrs {
		h.w.Write([]byte(fmt.Sprintf(" %s=%v", attr.k, attr.v)))
	}
	h.w.Write([]byte("\n"))
	return err
}

func (h *OtelHandler) With(attrs []Attr) Handler {
	m := map[string]any{}
	for k, v := range h.attrs {
		m[k] = v
	}
	for _, attr := range attrs {
		m[attr.k] = attr.v
	}

	return &OtelHandler{w: h.w, attrs: m}
}
