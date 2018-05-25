package simpletracer

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/satori/go.uuid"
)

const (
	traceHeaderName = "X-Trace-Id"
	spanHeaderName  = "X-Span-Id"
	CtxTracerName   = "CtxTracer"
)

type CtxTracer struct {
	TraceID string
	SpanID  string
}

type tracerHandler struct {
	handler http.Handler
}

func (c *tracerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	traceID := getTracerID(req)
	spanID := getSpanID(req)

	w.Header().Set(traceHeaderName, traceID)
	w.Header().Set(spanHeaderName, spanID)

	ctxTracer := &CtxTracer{
		TraceID: traceID,
		SpanID:  spanID,
	}
	ctx := context.WithValue(req.Context(), CtxTracerName, ctxTracer)
	c.handler.ServeHTTP(w, req.WithContext(ctx))
}

func GetCtxTrace(req *http.Request) *CtxTracer {
	if req != nil {
		ctx, ok := req.Context().Value("CtxTracerName").(*CtxTracer)
		if ok {
			return ctx
		}
	}
	return &CtxTracer{
		TraceID: getTracerID(nil),
		SpanID:  getSpanID(nil),
	}
}

func getTracerID(req *http.Request) string {
	if req != nil {
		traceID := req.Header.Get(traceHeaderName)
		if traceID != "" {
			return traceID
		}
	}

	uuID, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return uuID.String()
}

func getSpanID(req *http.Request) string {
	if req == nil {
		return "0"
	}
	spanID := req.Header.Get(spanHeaderName)
	i64, _ := strconv.ParseInt(spanID, 10, 64)
	return fmt.Sprintf("%d", i64+1)
}

func MiddlerWare(h http.Handler) http.Handler {
	handler := &tracerHandler{
		handler: h,
	}
	return handler
}
