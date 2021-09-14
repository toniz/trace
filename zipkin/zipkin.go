/*
 * Create By Xinwenjia 2021-09-11
 */

package zipkin

import (
    "net/http"
    "time"
    "context"
    "log"
    "os"
    "io/ioutil"
    "bytes"

    model "github.com/openzipkin/zipkin-go/model"
    zipkin "github.com/openzipkin/zipkin-go"
    zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"

    reporter "github.com/openzipkin/zipkin-go/reporter"
    logreporter "github.com/openzipkin/zipkin-go/reporter/log"
    httpreporter "github.com/openzipkin/zipkin-go/reporter/http"

    "github.com/bbthis/gosuit/trace"
)

const (
    ZipKinSpanKindUnspecified = string(model.Undetermined)
    ZipKinSpanKindInternal = string( model.Client)
    ZipKinSpanKindServer = string(model.Server)
    ZipKinSpanKindClient = string(model.Producer)
    ZipKinSpanKindProducer = string(model.Consumer)
);

var (
    TraceExportAlreadySet error = errors.New("TraceExportAlreadySet")
    TraceExportNotSupport error = errors.New("TraceExportNotSupport")
    TraceConfigCheckFailed error = errors.New("TraceConfigCheckFailed")
    TraceNotWork error = errors.New("TraceNotWork")
    TraceEndSpanFailed error = errors.New("TraceEndSpanFailed")
    TraceHttpCodeError error = errors.New("TraceHttpCodeError")
)

type ZipKin struct {
    tp *zipkin.Tracer
    report reporter.Reporter
    ZKNewHandler func(http.Handler) http.Handler
}

func init() {
    s := new(ZipKin)
    trace.SetTracer(s)
}

// Set GRPC Export
// ZipKin Not Support: github.com/openzipkin/zipkin-go
func (c *ZipKin) SetGrpcExport(ctx context.Context, filename, serviceName, version string) error {
    return TraceExportNotSupport
}

// Set Http Export
func (c *ZipKin) SetHttpExport(ctx context.Context, filename, serviceName, version string) error {
    if c.IsWork() == nil {
        return TraceExportAlreadySet
    }

    params := make(map[string]string)
    err := trace.LoadFile(filename, &params)
    if err != nil {
        return err
    }

    traceEndpoint := params["trace_endpoint"]
    ip := params["local_ip"]
    if len(traceEndpoint) == 0 || len(ip) == 0 {
        return TraceConfigCheckFailed
    }

    // Set Up A Span Http Reporter
    c.report = httpreporter.NewReporter(traceEndpoint)

    return c.SetGlobalProvider(ctx, serviceName, ip)
}

// Set Default Exprot: stdout
func (c *ZipKin) SetDefaultExport(ctx context.Context, serviceName, version string) error {
    if c.IsWork() == nil {
        return TraceExportAlreadySet
    }

    // Set Up A Span Log Reporter
    c.report = logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))

    return c.SetGlobalProvider(ctx, serviceName, version)
}

// Set Glocal Trace Provider
func (c *ZipKin) SetGlobalProvider(ctx context.Context, serviceName, ip string) error {
    // Create Our Local Service Endpoint
    endpoint, errNE := zipkin.NewEndpoint(serviceName, ip)
    if errNE != nil {
        return errNE
    }

    // Initialize Our Tracer
    var errNT error
    c.tp, errNT = zipkin.NewTracer(c.report, zipkin.WithLocalEndpoint(endpoint))
    if errNT != nil {
        return errNT
    }

    // Create Global Zipkin Http Server Middleware
    c.ZKNewHandler = zipkinhttp.NewServerMiddleware (
        c.tp, zipkinhttp.TagResponseSize(true),
    )

    return nil
}

func (c *ZipKin) Close(ctx context.Context) error {
    if c.report != nil {
        c.report.Close()
    }

    return nil
}

func (c *ZipKin) IsWork() error {
    if c.tp == nil {
        return TraceNotWork
    }

    return nil
}

func (c *ZipKin) NewSpan(ctx context.Context, name string, kind string) (context.Context, error) {
    if c.IsWork() != nil {
        return ctx, TraceNotWork
    }

    // Set Span Level To ctxc.
    level, _ := ctx.Value(trace.ContextSpanLevelKey).(int)
    ctx = context.WithValue(ctx, trace.ContextSpanLevelKey, level+1)

    _, ctxc := c.tp.StartSpanFromContext(ctx, name, zipkin.Kind(model.Kind(kind)))

    return ctxc, nil
}

// Call Span End
func (c *ZipKin) EndSpan(ctx context.Context) error {
    if c.IsWork() != nil {
        return TraceNotWork
    }

    span := zipkin.SpanFromContext(ctx)
    if span == nil {
        return TraceEndSpanFailed
    }

    span.Finish()
    return nil
}

// Adding An Attribute To A Span
func (c *ZipKin) AddSpanAttribute(ctx context.Context, params map[string]string) error {
    if c.IsWork() != nil {
        return TraceNotWork
    }

    span := zipkin.SpanFromContext(ctx)
    for k, v := range params {
        span.Tag(k, v)
    }

    return nil
}

// Adding An Event To A Span
func (c *ZipKin) AddSpanEvent(ctx context.Context, event string, params map[string]string) error {
    if c.IsWork() != nil {
        return TraceNotWork
    }

    span := zipkin.SpanFromContext(ctx)
    span.Annotate(time.Now(), event)
    for k, v := range params {
        span.Tag(k, v)
    }

    return nil
}

// ZipKin Not Support: github.com/openzipkin/zipkin-go
func (c *ZipKin) SetSpanOK(ctx context.Context, message string) error {
    return nil
}

// ZipKin Not Support: github.com/openzipkin/zipkin-go
func (c *ZipKin) SetSpanError(ctx context.Context, err error) error {
    return nil
}

// Http Handle Hook
func (c *ZipKin) NewHandler(handler http.Handler, name string) http.Handler {
    if c.IsWork() != nil {
        return handler
    }

    return c.ZKNewHandler(handler)
}

// Http Client Do Function
func (c *ZipKin) HttpDo(ctx context.Context, req *http.Request, name string, options ...interface{}) (*http.Response, error) {
    if c.IsWork() != nil {
        return nil, TraceNotWork
    }

    var opts []zipkinhttp.ClientOption
    for _, option := range options {
        if opt, ok := option.(zipkinhttp.ClientOption); ok {
            opts = append(opts, opt)
        }
    }

    client, errZK := zipkinhttp.NewClient(c.tp, opts...)
    if errZK != nil {
        return nil, errZK
    }

    res, errDO := client.DoWithAppSpan(req, name)
    if errDO != nil {
        return nil, errDO
    }

    return res, nil
}

func (c *ZipKin) HttpGet(ctx context.Context, url, name string) (string, error) {
    if c.IsWork() != nil {
        return "", TraceNotWork
    }

    req, errHttp := http.NewRequestWithContext(ctx, "GET", url, nil)
    if errHttp != nil {
        return "", errHttp
    }

    resp, err := c.HttpDo(ctx, req, name)
    if err != nil {
        return "", err
    }

    defer resp.Body.Close()
    resBody, errIO := ioutil.ReadAll(resp.Body)
    if errIO != nil {
        return "", errIO
    }

    if resp.StatusCode > 399 {
        return "", TraceHttpCodeError
    }

    return string(resBody), nil
}

func (c *ZipKin) HttpPost(ctx context.Context, url, contentType, body, name string) (string, error) {
    if c.IsWork() != nil {
        return "", TraceNotWork
    }

    reqBody := bytes.NewBuffer([]byte(body))
    req, errHttp := http.NewRequestWithContext(ctx, "POST", url, reqBody)
    if errHttp != nil {
        return "", errHttp
    }

    if len(contentType) != 0 {
        req.Header.Set("Content-Type", contentType)
    }

    resp, err := c.HttpDo(ctx, req, name)
    if err != nil {
        return "", err
    }

    defer resp.Body.Close()
    resBody, errIO := ioutil.ReadAll(resp.Body)
    if errIO != nil {
        return "", errIO
    }

    if resp.StatusCode > 399 {
        return "", TraceHttpCodeError
    }

    return string(resBody), nil
}
