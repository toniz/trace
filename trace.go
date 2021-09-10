/*
 * Create By Xinwenjia 2021-09-06
 */

package trace

import (
    "net/http"
    "context"
)

type Tracer interface {
    SetGrpcExport(ctx context.Context, filename, key, serviceName, version string) error
    SetDefaultExport(ctx context.Context, serviceName, version string) error
    Close(ctx context.Context) error
    IsWork() error
    NewHandler(handler http.Handler, name string) http.Handler
    NewSpan(ctx context.Context, name string, kind int) (context.Context, error)
    EndSpan(ctx context.Context) error
    AddSpanAttribute(ctx context.Context, params map[string]string) error
    AddSpanEvent(ctx context.Context, event string, params map[string]string) error
    SetSpanOK(ctx context.Context, message string) error
    SetSpanError(ctx context.Context, err error) error
}

var tracer Tracer

func SetTracer(l Tracer) {
    tracer = l
}

func SetGrpcExport(ctx context.Context, filename, key, serviceName, version string) error {
    return tracer.SetGrpcExport(ctx, filename, key, serviceName, version)
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

func NewHandler(handler http.Handler, name string) http.Handler {
    return tracer.NewHandler(handler, name)
}

func NewSpan(ctx context.Context, name string, kind int) (context.Context, error) {
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

