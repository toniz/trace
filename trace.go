/*
 * Create By Xinwenjia 2021-09-06
 */

package trace

import (
    "net/http"
    "context"
)

type Tracer interface {
    SetHttpExport(ctx context.Context, filename, serviceName, version string) error
    SetGrpcExport(ctx context.Context, filename, serviceName, version string) error
    SetDefaultExport(ctx context.Context, serviceName, version string) error
    Close(ctx context.Context) error
    IsWork() error

    NewSpan(ctx context.Context, name string, kind string) (context.Context, error)
    EndSpan(ctx context.Context) error
    AddSpanAttribute(ctx context.Context, params map[string]string) error
    AddSpanEvent(ctx context.Context, event string, params map[string]string) error
    SetSpanOK(ctx context.Context, message string) error
    SetSpanError(ctx context.Context, err error) error

    NewHandler(handler http.Handler, name string) http.Handler
    HttpDo(ctx context.Context, req *http.Request, name string, options ...interface{}) (*http.Response, error)
    HttpGet(ctx context.Context, url, name string) (string, error)
    HttpPost(ctx context.Context, url, contentType, body, name string) (string, error)
}

var tracer Tracer

const (
    ContextSpanLevelKey string = "span_level"
);

func SetTracer(l Tracer) {
    tracer = l
}

func SetHttpExport(ctx context.Context, filename, serviceName, version string) error {
    return tracer.SetHttpExport(ctx, filename, serviceName, version)
}

func SetGrpcExport(ctx context.Context, filename, serviceName, version string) error {
    return tracer.SetGrpcExport(ctx, filename, serviceName, version)
}

func SetDefaultExport(ctx context.Context, serviceName, version string) error {
    return tracer.SetDefaultExport(ctx, serviceName, version)
}

func Close(ctx context.Context) error {
    return tracer.Close(ctx)
}

func IsWork() error {
    return tracer.IsWork()
}

func NewSpan(ctx context.Context, name string, kind string) (context.Context, error) {
    return tracer.NewSpan(ctx, name, kind)
}

func EndSpan(ctx context.Context) error {
    return tracer.EndSpan(ctx)
}

func AddSpanAttribute(ctx context.Context, params map[string]string) error {
    return tracer.AddSpanAttribute(ctx, params)
}

func AddSpanEvent(ctx context.Context, event string, params map[string]string) error {
    return tracer.AddSpanEvent(ctx, event, params)
}

func SetSpanOK(ctx context.Context, message string) error {
    return tracer.SetSpanOK(ctx, message)
}

func SetSpanError(ctx context.Context, err error) error {
    return tracer.SetSpanError(ctx, err)
}

func NewHandler(handler http.Handler, name string) http.Handler {
    return tracer.NewHandler(handler, name)
}

func HttpDo(ctx context.Context, req *http.Request, name string, options ...interface{}) (*http.Response, error) {
    return tracer.HttpDo(ctx, req, name, options)
}

func HttpGet(ctx context.Context, url, name string) (string, error) {
    return tracer.HttpGet(ctx, url, name)
}

func HttpPost(ctx context.Context, url, contentType, body, name string) (string, error) {
    return tracer.HttpPost(ctx, url, contentType, body, name)
}


