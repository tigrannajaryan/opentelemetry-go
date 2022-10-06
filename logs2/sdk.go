package logs2

import (
	"context"
	"fmt"
	"io"

	"go.opentelemetry.io/otel/trace"
)

type OtelHandler struct {
	w     io.Writer
	attrs map[string]any
}

func NewOtelHandler(w io.Writer) *OtelHandler {
	return &OtelHandler{w: w}
}

func (h *OtelHandler) Handle(ctx context.Context, r Record) error {

	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		_, err := h.w.Write(
			[]byte(fmt.Sprintf(
				"%s trace_id=%v span_id=%v", r.msg,
				spanContext.TraceID(), spanContext.SpanID(),
			)),
		)
		if err != nil {
			return err
		}
	} else {
		h.w.Write([]byte(r.msg))
	}

	for k, v := range h.attrs {
		h.w.Write([]byte(fmt.Sprintf(" %s=%v", k, v)))
	}
	for _, attr := range r.attrs {
		h.w.Write([]byte(fmt.Sprintf(" %s=%v", attr.k, attr.v)))
	}
	h.w.Write([]byte("\n"))
	return nil
}

// With returns a new Handler whose attributes consist of
// the receiver's attributes concatenated with the arguments.
func (h *OtelHandler) With(attrs []Attr) Handler {
	m := map[string]any{}
	for k, v := range h.attrs {
		m[k] = v
	}
	for _, attr := range attrs {
		m[attr.k] = attr.v
	}

	return &TextHandler{w: h.w, attrs: m}
}
